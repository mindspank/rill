package transporter

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/rilldata/rill/runtime/drivers"
	"go.uber.org/zap"
)

type motherduckToDuckDB struct {
	to     drivers.OLAPStore
	from   drivers.Handle
	logger *zap.Logger
}

var _ drivers.Transporter = &motherduckToDuckDB{}

func NewMotherduckToDuckDB(from drivers.Handle, to drivers.OLAPStore, logger *zap.Logger) drivers.Transporter {
	return &motherduckToDuckDB{
		to:     to,
		from:   from,
		logger: logger,
	}
}

// TODO: should it run count from user_query to set target in progress ?
func (t *motherduckToDuckDB) Transfer(ctx context.Context, srcProps, sinkProps map[string]any, opts *drivers.TransferOptions) error {
	srcCfg, err := parseDBSourceProperties(srcProps)
	if err != nil {
		return err
	}

	sinkCfg, err := parseSinkProperties(sinkProps)
	if err != nil {
		return err
	}

	config := t.from.Config()
	err = t.to.WithConnection(ctx, 1, true, false, func(ctx, ensuredCtx context.Context, _ *sql.Conn) error {
		res, err := t.to.Execute(ctx, &drivers.Statement{Query: "SELECT current_database();"})
		if err != nil {
			return err
		}
		defer res.Close()

		res.Next()
		var localDB string
		if err := res.Scan(&localDB); err != nil {
			return err
		}

		// get token
		token := config["token"]
		if token == "" && config["allow_host_access"].(bool) {
			token = os.Getenv("motherduck_token")
		}
		if token == "" {
			return fmt.Errorf("no motherduck token found. Refer to this documentation for instructions: https://docs.rilldata.com/deploy/credentials/motherduck")
		}

		// load motherduck extension; connect to motherduck service
		err = t.to.Exec(ctx, &drivers.Statement{Query: "INSTALL 'motherduck'; LOAD 'motherduck';"})
		if err != nil {
			return fmt.Errorf("failed to load motherduck extension %w", err)
		}

		err = t.to.Exec(ctx, &drivers.Statement{Query: fmt.Sprintf("PRAGMA MD_CONNECT('token=%s');", token)})
		if err != nil {
			if !strings.Contains(err.Error(), "already connected") {
				return fmt.Errorf("failed to connect to motherduck %w", err)
			}
		}

		var names []string

		db := srcCfg.Database
		if db == "" {
			// get list of all motherduck databases
			res, err = t.to.Execute(ctx, &drivers.Statement{Query: "SELECT name FROM md_databases();"})
			if err != nil {
				return err
			}
			defer res.Close()

			for res.Next() {
				var name string
				if res.Scan(&name) != nil {
					return err
				}
				names = append(names, name)
			}
			// single motherduck db, use db to allow user to run query without specifying db name
			if len(names) == 1 {
				db = names[0]
			}
		}

		if db != "" {
			err = t.to.Exec(ctx, &drivers.Statement{Query: fmt.Sprintf("USE %s;", safeName(db))})
			if err != nil {
				return err
			}

			defer func(ctx context.Context) { // revert back to localdb
				err = t.to.Exec(ctx, &drivers.Statement{Query: fmt.Sprintf("USE %s;", safeName(localDB))})
				if err != nil {
					t.logger.Error("failed to switch to local database", zap.Error(err))
				}
			}(ensuredCtx)
		}

		if srcCfg.SQL == "" {
			return fmt.Errorf("property \"sql\" is mandatory for connector \"motherduck\"")
		}

		userQuery := strings.TrimSpace(srcCfg.SQL)
		userQuery, _ = strings.CutSuffix(userQuery, ";") // trim trailing semi colon
		query := fmt.Sprintf("CREATE OR REPLACE TABLE %s.%s AS (%s);", safeName(localDB), safeName(sinkCfg.Table), userQuery)
		return t.to.Exec(ctx, &drivers.Statement{Query: query})
	})
	return err
}
