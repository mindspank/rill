package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/rilldata/rill/runtime/api"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/services/catalog"
	"github.com/rilldata/rill/runtime/services/catalog/artifacts"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func CreateSource(t *testing.T, s *catalog.Service, name string, file string, path string) string {
	ctx := context.Background()
	err := artifacts.Write(ctx, s.Repo, s.RepoId, &api.CatalogObject{
		Name: name,
		Type: api.CatalogObject_TYPE_SOURCE,
		Source: &api.Source{
			Name:      name,
			Connector: "file",
			Properties: toProtoStruct(map[string]any{
				"path": file,
			}),
		},
		Path: path,
	})
	require.NoError(t, err)
	blob, err := s.Repo.Get(ctx, s.RepoId, path)
	require.NoError(t, err)
	return blob
}

func CreateModel(t *testing.T, s *catalog.Service, name string, sql string, path string) string {
	ctx := context.Background()
	err := artifacts.Write(ctx, s.Repo, s.RepoId, &api.CatalogObject{
		Name: name,
		Type: api.CatalogObject_TYPE_MODEL,
		Model: &api.Model{
			Name:    name,
			Sql:     sql,
			Dialect: api.Model_DIALECT_DUCKDB,
		},
		Path: path,
	})
	require.NoError(t, err)
	blob, err := s.Repo.Get(ctx, s.RepoId, path)
	require.NoError(t, err)
	return blob
}

func CreateMetricsView(t *testing.T, s *catalog.Service, metricsView *api.MetricsView, path string) string {
	ctx := context.Background()
	err := artifacts.Write(ctx, s.Repo, s.RepoId, &api.CatalogObject{
		Name:        metricsView.Name,
		Type:        api.CatalogObject_TYPE_METRICS_VIEW,
		MetricsView: metricsView,
		Path:        path,
	})
	require.NoError(t, err)
	blob, err := s.Repo.Get(ctx, s.RepoId, path)
	require.NoError(t, err)
	return blob
}

func toProtoStruct(obj map[string]any) *structpb.Struct {
	s, err := structpb.NewStruct(obj)
	if err != nil {
		panic(err)
	}
	return s
}

func AssertTable(t *testing.T, s *catalog.Service, name string, path string) {
	AssertInCatalogStore(t, s, name, path)

	rows, err := s.Olap.Execute(context.Background(), &drivers.Statement{
		Query:    fmt.Sprintf("select count(*) as count from %s", name),
		Args:     nil,
		DryRun:   false,
		Priority: 0,
	})
	require.NoError(t, err)
	var count int
	rows.Next()
	require.NoError(t, rows.Scan(&count))
	require.Greater(t, count, 1)
	require.NoError(t, rows.Close())

	table, err := s.Olap.InformationSchema().Lookup(context.Background(), name)
	require.NoError(t, err)
	require.Equal(t, name, table.Name)
}

func AssertInCatalogStore(t *testing.T, s *catalog.Service, name string, path string) {
	catalogObject, ok := s.Catalog.FindObject(context.Background(), s.InstId, name)
	require.True(t, ok)
	require.Equal(t, name, catalogObject.Name)
	require.Equal(t, path, catalogObject.Path)
}

func AssertTableAbsence(t *testing.T, s *catalog.Service, name string) {
	_, ok := s.Catalog.FindObject(context.Background(), s.InstId, name)
	require.False(t, ok)

	_, err := s.Olap.InformationSchema().Lookup(context.Background(), name)
	require.ErrorIs(t, err, drivers.ErrNotFound)
}

func AssertMigration(t *testing.T, result *catalog.MigrationResult, errCount int, addCount int, updateCount int, dropCount int) {
	require.Len(t, result.Errors, errCount)
	require.Len(t, result.AddedObjects, addCount)
	require.Len(t, result.UpdatedObjects, updateCount)
	require.Len(t, result.DroppedObjects, dropCount)
}
