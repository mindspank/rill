package queries

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetricsViewAggregation struct {
	MetricsViewName    string                                       `json:"metrics_view,omitempty"`
	Dimensions         []*runtimev1.MetricsViewAggregationDimension `json:"dimensions,omitempty"`
	Measures           []*runtimev1.MetricsViewAggregationMeasure   `json:"measures,omitempty"`
	Sort               []*runtimev1.MetricsViewAggregationSort      `json:"sort,omitempty"`
	TimeStart          *timestamppb.Timestamp                       `json:"time_start,omitempty"`
	TimeEnd            *timestamppb.Timestamp                       `json:"time_end,omitempty"`
	Filter             *runtimev1.MetricsViewFilter                 `json:"filter,omitempty"`
	Priority           int32                                        `json:"priority,omitempty"`
	Limit              *int64                                       `json:"limit,omitempty"`
	Offset             int64                                        `json:"offset,omitempty"`
	MetricsView        *runtimev1.MetricsView                       `json:"-"`
	ResolvedMVSecurity *runtime.ResolvedMetricsViewSecurity         `json:"security"`

	Result *runtimev1.MetricsViewAggregationResponse `json:"-"`
}

var _ runtime.Query = &MetricsViewAggregation{}

func (q *MetricsViewAggregation) Key() string {
	r, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("MetricsViewAggregation:%s", string(r))
}

func (q *MetricsViewAggregation) Deps() []string {
	return []string{q.MetricsViewName}
}

func (q *MetricsViewAggregation) MarshalResult() *runtime.QueryResult {
	return &runtime.QueryResult{
		Value: q.Result,
		Bytes: sizeProtoMessage(q.Result),
	}
}

func (q *MetricsViewAggregation) UnmarshalResult(v any) error {
	res, ok := v.(*runtimev1.MetricsViewAggregationResponse)
	if !ok {
		return fmt.Errorf("MetricsViewAggregation: mismatched unmarshal input")
	}
	q.Result = res
	return nil
}

func (q *MetricsViewAggregation) Resolve(ctx context.Context, rt *runtime.Runtime, instanceID string, priority int) error {
	olap, release, err := rt.OLAP(ctx, instanceID)
	if err != nil {
		return err
	}
	defer release()

	if olap.Dialect() != drivers.DialectDuckDB && olap.Dialect() != drivers.DialectDruid {
		return fmt.Errorf("not available for dialect '%s'", olap.Dialect())
	}

	if q.MetricsView.TimeDimension == "" && (q.TimeStart != nil || q.TimeEnd != nil) {
		return fmt.Errorf("metrics view '%s' does not have a time dimension", q.MetricsView)
	}

	// Build query
	sql, args, err := q.buildMetricsAggregationSQL(q.MetricsView, olap.Dialect(), q.ResolvedMVSecurity)
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	// Execute
	schema, data, err := olapQuery(ctx, olap, priority, sql, args)
	if err != nil {
		return err
	}

	q.Result = &runtimev1.MetricsViewAggregationResponse{
		Schema: schema,
		Data:   data,
	}

	return nil
}

func (q *MetricsViewAggregation) Export(ctx context.Context, rt *runtime.Runtime, instanceID string, w io.Writer, opts *runtime.ExportOptions) error {
	err := q.Resolve(ctx, rt, instanceID, opts.Priority)
	if err != nil {
		return err
	}

	filename := strings.ReplaceAll(q.MetricsView.Model, `"`, `_`)
	if q.TimeStart != nil || q.TimeEnd != nil || q.Filter != nil && (len(q.Filter.Include) > 0 || len(q.Filter.Exclude) > 0) {
		filename += "_filtered"
	}

	meta := structTypeToMetricsViewColumn(q.Result.Schema)

	if opts.PreWriteHook != nil {
		err = opts.PreWriteHook(filename)
		if err != nil {
			return err
		}
	}

	switch opts.Format {
	case runtimev1.ExportFormat_EXPORT_FORMAT_UNSPECIFIED:
		return fmt.Errorf("unspecified format")
	case runtimev1.ExportFormat_EXPORT_FORMAT_CSV:
		return writeCSV(meta, q.Result.Data, w)
	case runtimev1.ExportFormat_EXPORT_FORMAT_XLSX:
		return writeXLSX(meta, q.Result.Data, w)
	case runtimev1.ExportFormat_EXPORT_FORMAT_PARQUET:
		return writeParquet(meta, q.Result.Data, w)
	}

	return nil
}

func (q *MetricsViewAggregation) buildMetricsAggregationSQL(mv *runtimev1.MetricsView, dialect drivers.Dialect, policy *runtime.ResolvedMetricsViewSecurity) (string, []any, error) {
	if len(q.Dimensions) == 0 && len(q.Measures) == 0 {
		return "", nil, errors.New("no dimensions or measures specified")
	}

	selectCols := make([]string, 0, len(q.Dimensions)+len(q.Measures))
	groupCols := make([]string, 0, len(q.Dimensions))
	args := []any{}

	for _, d := range q.Dimensions {
		// Handle regular dimensions
		if d.TimeGrain == runtimev1.TimeGrain_TIME_GRAIN_UNSPECIFIED {
			col, err := metricsViewDimensionToSafeColumn(mv, d.Name)
			if err != nil {
				return "", nil, err
			}

			selectCols = append(selectCols, fmt.Sprintf("%s as %s", col, safeName(d.Name)))
			groupCols = append(groupCols, col)
			continue
		}

		// Handle time dimension
		expr, exprArgs, err := q.buildTimestampExpr(d, dialect)
		if err != nil {
			return "", nil, err
		}
		selectCols = append(selectCols, fmt.Sprintf("%s as %s", expr, safeName(d.Name)))
		groupCols = append(groupCols, expr)
		args = append(args, exprArgs...)
	}

	for _, m := range q.Measures {
		switch m.BuiltinMeasure {
		case runtimev1.BuiltinMeasure_BUILTIN_MEASURE_UNSPECIFIED:
			expr, err := metricsViewMeasureExpression(mv, m.Name)
			if err != nil {
				return "", nil, err
			}
			selectCols = append(selectCols, fmt.Sprintf("%s as %s", expr, safeName(m.Name)))
		case runtimev1.BuiltinMeasure_BUILTIN_MEASURE_COUNT:
			selectCols = append(selectCols, fmt.Sprintf("COUNT(*) as %s", safeName(m.Name)))
		case runtimev1.BuiltinMeasure_BUILTIN_MEASURE_COUNT_DISTINCT:
			if len(m.BuiltinMeasureArgs) != 1 {
				return "", nil, fmt.Errorf("builtin measure '%s' expects 1 argument", m.BuiltinMeasure.String())
			}
			arg := m.BuiltinMeasureArgs[0].GetStringValue()
			if arg == "" {
				return "", nil, fmt.Errorf("builtin measure '%s' expects non-empty string argument, got '%v'", m.BuiltinMeasure.String(), m.BuiltinMeasureArgs[0])
			}
			selectCols = append(selectCols, fmt.Sprintf("COUNT(DISTINCT %s) as %s", safeName(arg), safeName(m.Name)))
		default:
			return "", nil, fmt.Errorf("unknown builtin measure '%d'", m.BuiltinMeasure)
		}
	}

	groupClause := ""
	if len(groupCols) > 0 {
		groupClause = "GROUP BY " + strings.Join(groupCols, ", ")
	}

	whereClause := ""
	if mv.TimeDimension != "" {
		if q.TimeStart != nil {
			whereClause += fmt.Sprintf(" AND %s >= ?", safeName(mv.TimeDimension))
			args = append(args, q.TimeStart.AsTime())
		}
		if q.TimeEnd != nil {
			whereClause += fmt.Sprintf(" AND %s < ?", safeName(mv.TimeDimension))
			args = append(args, q.TimeEnd.AsTime())
		}
	}
	if q.Filter != nil {
		clause, clauseArgs, err := buildFilterClauseForMetricsViewFilter(mv, q.Filter, dialect, policy)
		if err != nil {
			return "", nil, err
		}
		whereClause += " " + clause
		args = append(args, clauseArgs...)
	}
	if len(whereClause) > 0 {
		whereClause = "WHERE 1=1" + whereClause
	}

	sortingCriteria := make([]string, 0, len(q.Sort))
	for _, s := range q.Sort {
		sortCriterion := safeName(s.Name)
		if s.Desc {
			sortCriterion += " DESC"
		}
		if dialect == drivers.DialectDuckDB {
			sortCriterion += " NULLS LAST"
		}
		sortingCriteria = append(sortingCriteria, sortCriterion)
	}
	orderClause := ""
	if len(sortingCriteria) > 0 {
		orderClause = "ORDER BY " + strings.Join(sortingCriteria, ", ")
	}

	var limitClause string
	if q.Limit != nil {
		if *q.Limit == 0 {
			*q.Limit = 100
		}
		limitClause = fmt.Sprintf("LIMIT %d", *q.Limit)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s %s %s %s %s OFFSET %d",
		strings.Join(selectCols, ", "),
		safeName(mv.Model),
		whereClause,
		groupClause,
		orderClause,
		limitClause,
		q.Offset,
	)

	return sql, args, nil
}

func (q *MetricsViewAggregation) buildTimestampExpr(dim *runtimev1.MetricsViewAggregationDimension, dialect drivers.Dialect) (string, []any, error) {
	var colName string
	if dim.Name == q.MetricsView.TimeDimension {
		colName = dim.Name
	} else {
		col, err := metricsViewDimensionToSafeColumn(q.MetricsView, dim.Name)
		if err != nil {
			return "", nil, err
		}
		colName = col
	}

	switch dialect {
	case drivers.DialectDuckDB:
		if dim.TimeZone == "" || dim.TimeZone == "UTC" {
			return fmt.Sprintf("date_trunc('%s', %s)", convertToDateTruncSpecifier(dim.TimeGrain), safeName(colName)), nil, nil
		}
		return fmt.Sprintf("timezone(?, date_trunc('%s', timezone(?, %s::TIMESTAMPTZ)))", convertToDateTruncSpecifier(dim.TimeGrain), safeName(colName)), []any{dim.TimeZone, dim.TimeZone}, nil
	case drivers.DialectDruid:
		if dim.TimeZone == "" || dim.TimeZone == "UTC" {
			return fmt.Sprintf("date_trunc('%s', %s)", convertToDateTruncSpecifier(dim.TimeGrain), safeName(colName)), nil, nil
		}
		return fmt.Sprintf("time_floor(%s, '%s', null, CAST(? AS VARCHAR)))", safeName(colName), convertToDruidTimeFloorSpecifier(dim.TimeGrain)), []any{dim.TimeZone}, nil
	default:
		return "", nil, fmt.Errorf("unsupported dialect %q", dialect)
	}
}
