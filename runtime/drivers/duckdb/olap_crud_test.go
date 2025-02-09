package duckdb

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_connection_CreateTableAsSelect(t *testing.T) {
	temp := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(temp, "default"), fs.ModePerm))
	dbPath := filepath.Join(temp, "default", "normal.db")
	handle, err := Driver{}.Open(map[string]any{"dsn": dbPath}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	normalConn := handle.(*connection)
	normalConn.AsOLAP("default")
	require.NoError(t, normalConn.Migrate(context.Background()))

	dbPath = filepath.Join(temp, "default", "view.db")
	handle, err = Driver{}.Open(map[string]any{"dsn": dbPath, "external_table_storage": true}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	viewConnection := handle.(*connection)
	require.NoError(t, viewConnection.Migrate(context.Background()))
	viewConnection.AsOLAP("default")

	tests := []struct {
		testName    string
		name        string
		view        bool
		tableAsView bool
		c           *connection
	}{
		{
			testName: "select from view",
			name:     "my-view",
			view:     true,
			c:        normalConn,
		},
		{
			testName: "select from table",
			name:     "my-table",
			c:        normalConn,
		},
		{
			testName: "select from view with external_table_storage",
			name:     "my-view",
			c:        viewConnection,
			view:     true,
		},
		{
			testName:    "select from table with external_table_storage",
			name:        "my-table",
			c:           viewConnection,
			tableAsView: true,
		},
	}
	ctx := context.Background()
	sql := "SELECT 1"
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := tt.c.CreateTableAsSelect(ctx, tt.name, tt.view, sql)
			require.NoError(t, err)
			res, err := tt.c.Execute(ctx, &drivers.Statement{Query: fmt.Sprintf("SELECT count(*) FROM %q", tt.name)})
			require.NoError(t, err)
			require.True(t, res.Next())
			var count int
			require.NoError(t, res.Scan(&count))
			require.Equal(t, 1, count)
			require.NoError(t, res.Close())

			if tt.tableAsView {
				res, err := tt.c.Execute(ctx, &drivers.Statement{Query: fmt.Sprintf("SELECT count(*) FROM information_schema.tables WHERE table_name='%s' AND table_type='VIEW'", tt.name)})
				require.NoError(t, err)
				require.True(t, res.Next())
				var count int
				require.NoError(t, res.Scan(&count))
				require.Equal(t, 1, count)
				require.NoError(t, res.Close())
				contents, err := os.ReadFile(filepath.Join(temp, "default", tt.name, "version.txt"))
				require.NoError(t, err)
				version, err := strconv.ParseInt(string(contents), 10, 64)
				require.NoError(t, err)
				require.True(t, time.Since(time.UnixMilli(version)).Seconds() < 1)
			}
		})
	}
}

func Test_connection_CreateTableAsSelectMultipleTimes(t *testing.T) {
	temp := t.TempDir()

	dbPath := filepath.Join(temp, "view.db")
	handle, err := Driver{}.Open(map[string]any{"dsn": dbPath, "external_table_storage": true}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	c := handle.(*connection)
	require.NoError(t, c.Migrate(context.Background()))
	c.AsOLAP("default")

	err = c.CreateTableAsSelect(context.Background(), "test-select-multiple", false, "select 1")
	require.NoError(t, err)
	time.Sleep(2 * time.Millisecond)
	err = c.CreateTableAsSelect(context.Background(), "test-select-multiple", false, "select 'hello'")
	require.NoError(t, err)

	dirs, err := os.ReadDir(filepath.Join(temp, "test-select-multiple"))
	require.NoError(t, err)
	names := make([]string, 0)
	for _, dir := range dirs {
		names = append(names, dir.Name())
	}

	err = c.CreateTableAsSelect(context.Background(), "test-select-multiple", false, "select fail query")
	require.Error(t, err)

	dirs, err = os.ReadDir(filepath.Join(temp, "test-select-multiple"))
	require.NoError(t, err)
	newNames := make([]string, 0)
	for _, dir := range dirs {
		newNames = append(newNames, dir.Name())
	}

	require.Equal(t, names, newNames)

	res, err := c.Execute(context.Background(), &drivers.Statement{Query: fmt.Sprintf("SELECT * FROM %q", "test-select-multiple")})
	require.NoError(t, err)
	require.True(t, res.Next())
	var name string
	require.NoError(t, res.Scan(&name))
	require.Equal(t, "hello", name)
	require.False(t, res.Next())
	require.NoError(t, res.Close())
}

func Test_connection_DropTable(t *testing.T) {
	temp := t.TempDir()

	dbPath := filepath.Join(temp, "view.db")
	handle, err := Driver{}.Open(map[string]any{"dsn": dbPath, "external_table_storage": true}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	c := handle.(*connection)
	require.NoError(t, c.Migrate(context.Background()))
	c.AsOLAP("default")

	err = c.CreateTableAsSelect(context.Background(), "test-drop", false, "select 1")
	require.NoError(t, err)

	err = c.DropTable(context.Background(), "test-drop", false)
	require.NoError(t, err)

	_, err = os.ReadDir(filepath.Join(temp, "test-drop"))
	require.True(t, os.IsNotExist(err))

	res, err := c.Execute(context.Background(), &drivers.Statement{Query: "SELECT count(*) FROM information_schema.tables WHERE table_name='test-drop' AND table_type='VIEW'"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var count int
	require.NoError(t, res.Scan(&count))
	require.Equal(t, 0, count)
	require.NoError(t, res.Close())
}

func Test_connection_InsertTableAsSelect(t *testing.T) {
	temp := t.TempDir()

	dbPath := filepath.Join(temp, "view.db")
	handle, err := Driver{}.Open(map[string]any{"dsn": dbPath, "external_table_storage": true}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	c := handle.(*connection)
	require.NoError(t, c.Migrate(context.Background()))
	c.AsOLAP("default")

	err = c.CreateTableAsSelect(context.Background(), "test-insert", false, "select 1")
	require.NoError(t, err)

	err = c.InsertTableAsSelect(context.Background(), "test-insert", false, "select 2")
	require.NoError(t, err)

	err = c.InsertTableAsSelect(context.Background(), "test-insert", true, "select 3")
	require.Error(t, err)

	res, err := c.Execute(context.Background(), &drivers.Statement{Query: "SELECT count(*) FROM 'test-insert'"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var count int
	require.NoError(t, res.Scan(&count))
	require.Equal(t, 2, count)
	require.NoError(t, res.Close())
}

func Test_connection_RenameTable(t *testing.T) {
	temp := t.TempDir()
	os.Mkdir(temp, fs.ModePerm)

	dbPath := filepath.Join(temp, "view.db")
	handle, err := Driver{}.Open(map[string]any{"dsn": dbPath, "external_table_storage": true}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	c := handle.(*connection)
	require.NoError(t, c.Migrate(context.Background()))
	c.AsOLAP("default")

	err = c.CreateTableAsSelect(context.Background(), "test-rename", false, "select 1")
	require.NoError(t, err)

	err = c.RenameTable(context.Background(), "test-rename", "rename-test", false)
	require.NoError(t, err)

	res, err := c.Execute(context.Background(), &drivers.Statement{Query: "SELECT count(*) FROM 'rename-test'"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var count int
	require.NoError(t, res.Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, res.Close())
}

func Test_connection_RenameToExistingTable(t *testing.T) {
	temp := t.TempDir()
	os.Mkdir(temp, fs.ModePerm)

	dbPath := filepath.Join(temp, "default", "view.db")
	handle, err := Driver{}.Open(map[string]any{"dsn": dbPath, "external_table_storage": true}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	c := handle.(*connection)
	require.NoError(t, c.Migrate(context.Background()))
	c.AsOLAP("default")

	err = c.CreateTableAsSelect(context.Background(), "source", false, "SELECT 1 AS data")
	require.NoError(t, err)

	err = c.CreateTableAsSelect(context.Background(), "_tmp_source", false, "SELECT 2 AS DATA")
	require.NoError(t, err)

	err = c.RenameTable(context.Background(), "_tmp_source", "source", false)
	require.NoError(t, err)

	res, err := c.Execute(context.Background(), &drivers.Statement{Query: "SELECT * FROM 'source'"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var num int
	require.NoError(t, res.Scan(&num))
	require.Equal(t, 2, num)
	require.NoError(t, res.Close())
}

func Test_connection_AddTableColumn(t *testing.T) {
	temp := t.TempDir()
	os.Mkdir(temp, fs.ModePerm)

	dbPath := filepath.Join(temp, "view.db")
	handle, err := Driver{}.Open(map[string]any{"dsn": dbPath, "external_table_storage": true}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	c := handle.(*connection)
	require.NoError(t, c.Migrate(context.Background()))
	c.AsOLAP("default")

	err = c.CreateTableAsSelect(context.Background(), "test alter column", false, "select 1 as data")
	require.NoError(t, err)

	res, err := c.Execute(context.Background(), &drivers.Statement{Query: "SELECT data_type FROM information_schema.columns WHERE table_name='test alter column' AND table_catalog = 'view'"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var typ string
	require.NoError(t, res.Scan(&typ))
	require.Equal(t, "INTEGER", typ)
	require.NoError(t, res.Close())

	err = c.AlterTableColumn(context.Background(), "test alter column", "data", "VARCHAR")
	require.NoError(t, err)

	res, err = c.Execute(context.Background(), &drivers.Statement{Query: "SELECT data_type FROM information_schema.columns WHERE table_name='test alter column' AND table_catalog = 'view'"})
	require.NoError(t, err)
	require.True(t, res.Next())
	require.NoError(t, res.Scan(&typ))
	require.Equal(t, "VARCHAR", typ)
	require.NoError(t, res.Close())
}

func Test_connection_RenameToExistingTableOld(t *testing.T) {
	handle, err := Driver{}.Open(map[string]any{"dsn": ""}, false, activity.NewNoopClient(), zap.NewNop())
	require.NoError(t, err)
	c := handle.(*connection)
	require.NoError(t, c.Migrate(context.Background()))
	c.AsOLAP("default")

	err = c.CreateTableAsSelect(context.Background(), "source", false, "SELECT 1 AS data")
	require.NoError(t, err)

	err = c.CreateTableAsSelect(context.Background(), "_tmp_source", false, "SELECT 2 AS DATA")
	require.NoError(t, err)

	err = c.RenameTable(context.Background(), "_tmp_source", "source", false)
	require.NoError(t, err)

	res, err := c.Execute(context.Background(), &drivers.Statement{Query: "SELECT * FROM 'source'"})
	require.NoError(t, err)
	require.True(t, res.Next())
	var num int
	require.NoError(t, res.Scan(&num))
	require.Equal(t, 2, num)
	require.NoError(t, res.Close())
}
