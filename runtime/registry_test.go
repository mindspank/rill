package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/c2h5oh/datasize"
	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/activity"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

func TestRuntime_EditInstance(t *testing.T) {
	repodsn := t.TempDir()
	newRepodsn := t.TempDir()
	rt := NewTestRunTime(t)
	tests := []struct {
		name      string
		inst      *drivers.Instance
		wantErr   bool
		savedInst *drivers.Instance
	}{
		{
			name: "edit env",
			inst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  true,
				Variables:     map[string]string{"connectors.s3.region": "us-east-1"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
			savedInst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  true,
				Variables:     map[string]string{"connectors.s3.region": "us-east-1", "allow_host_access": "true"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
		},
		{
			name: "edit drivers",
			inst: &drivers.Instance{
				OLAPConnector: "olap1",
				RepoConnector: "repo1",
				EmbedCatalog:  true,
				Variables:     map[string]string{"host": "localhost"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo1",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "olap1",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
			savedInst: &drivers.Instance{
				OLAPConnector: "olap1",
				RepoConnector: "repo1",
				EmbedCatalog:  true,
				Variables:     map[string]string{"host": "localhost", "allow_host_access": "true"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo1",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "olap1",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
		},
		{
			name: "edit env and embed catalog",
			inst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"connectors.s3.region": "us-east-1"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
			savedInst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"connectors.s3.region": "us-east-1", "allow_host_access": "true"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
		},
		{
			name: "edit olap dsn",
			inst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"host": "localhost"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": "?access_mode=read_write"},
					},
				},
			},
			savedInst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"host": "localhost", "allow_host_access": "true"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": "?access_mode=read_write"},
					},
				},
			},
		},
		{
			name: "edit repo dsn",
			inst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"host": "localhost"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": newRepodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
			savedInst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"host": "localhost", "allow_host_access": "true"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": newRepodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
			},
		},
		{
			name: "edit annotations",
			inst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"host": "localhost"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": newRepodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
				Annotations: map[string]string{
					"organization_name": "org_name",
				},
			},
			savedInst: &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  false,
				Variables:     map[string]string{"host": "localhost", "allow_host_access": "true"},
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": newRepodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
				Annotations: map[string]string{
					"organization_name": "org_name",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			//create instance
			inst := &drivers.Instance{
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  true,
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": ""},
					},
				},
			}
			require.NoError(t, rt.CreateInstance(context.Background(), inst))
			// load all caches
			svc, err := rt.NewCatalogService(ctx, inst.ID)
			require.NoError(t, err)

			// edit instance
			tt.inst.ID = inst.ID
			err = rt.EditInstance(ctx, tt.inst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Runtime.EditInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			// verify db instances are correctly updated
			newInst, err := rt.FindInstance(ctx, inst.ID)
			require.NoError(t, err)
			require.Equal(t, inst.ID, newInst.ID)
			require.Equal(t, tt.savedInst.OLAPConnector, newInst.OLAPConnector)
			require.True(t, connectorsEqual(tt.savedInst.Connectors[0], newInst.Connectors[0]))
			require.True(t, connectorsEqual(tt.savedInst.Connectors[1], newInst.Connectors[1]))
			require.Equal(t, tt.savedInst.RepoConnector, newInst.RepoConnector)
			require.Equal(t, tt.savedInst.EmbedCatalog, newInst.EmbedCatalog)
			require.Greater(t, time.Since(newInst.CreatedOn), time.Since(newInst.UpdatedOn))
			require.Equal(t, tt.savedInst.Variables, newInst.Variables)

			// verify older olap connection is closed
			driver, cfg, err := rt.connectorConfig(ctx, inst.ID, inst.OLAPConnector)
			_, ok := rt.connCache.acquired[inst.ID+driver+generateKey(cfg)]
			require.Equal(t, false, ok)
			driver, cfg, err = rt.connectorConfig(ctx, inst.ID, inst.RepoConnector)
			_, ok = rt.connCache.acquired[inst.ID+driver+generateKey(cfg)]
			require.Equal(t, false, ok)
			_, ok = rt.migrationMetaCache.cache.Get(inst.ID)
			require.Equal(t, false, ok)
			_, err = svc.Olap.Execute(context.Background(), &drivers.Statement{Query: "SELECT COUNT(*) FROM rill.migration_version"})
			require.Equal(t, true, err != nil)
		})
	}
}

func TestRuntime_DeleteInstance(t *testing.T) {
	repodsn := t.TempDir()
	rt := NewTestRunTime(t)
	tests := []struct {
		name       string
		instanceID string
		dropDB     bool
		wantErr    bool
	}{
		{"delete valid no drop", "default", false, false},
		{"delete valid drop", "default", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create test data
			ctx := context.Background()
			dbFile := filepath.Join(t.TempDir(), "test.db")
			inst := &drivers.Instance{
				ID:            "default",
				OLAPConnector: "duckdb",
				RepoConnector: "repo",
				EmbedCatalog:  true,
				Connectors: []*runtimev1.Connector{
					{
						Type:   "file",
						Name:   "repo",
						Config: map[string]string{"dsn": repodsn},
					},
					{
						Type:   "duckdb",
						Name:   "duckdb",
						Config: map[string]string{"dsn": dbFile},
					},
				},
			}
			require.NoError(t, rt.CreateInstance(context.Background(), inst))
			// load all caches
			svc, err := rt.NewCatalogService(ctx, inst.ID)
			require.NoError(t, err)

			// ingest some data
			require.NoError(t, svc.Olap.Exec(ctx, &drivers.Statement{Query: "CREATE TABLE data(id INTEGER, name VARCHAR)"}))
			require.NoError(t, svc.Olap.Exec(ctx, &drivers.Statement{Query: "INSERT INTO data VALUES (1, 'Mark'), (2, 'Hannes')"}))
			require.NoError(t, svc.Catalog.CreateEntry(ctx, &drivers.CatalogEntry{
				Name: "data",
				Type: drivers.ObjectTypeTable,
				Object: &runtimev1.Table{
					Name:    "data",
					Managed: true,
				},
			}))
			require.ErrorContains(t, svc.Catalog.CreateEntry(ctx, &drivers.CatalogEntry{
				Name: "data",
				Type: drivers.ObjectTypeModel,
				Object: &runtimev1.Model{
					Name: "data",
				},
			}), "catalog entry with name \"data\" already exists")

			// delete instance
			err = rt.DeleteInstance(ctx, tt.instanceID, tt.dropDB)
			if (err != nil) != tt.wantErr {
				t.Errorf("Runtime.DeleteInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			// verify db is correctly cleared
			_, err = rt.FindInstance(ctx, inst.ID)
			require.Error(t, err)

			// verify older olap connection is closed and cache updated
			require.False(t, rt.connCache.lru.Contains(inst.ID+"duckdb"+fmt.Sprintf("dsn:%s ", dbFile)))
			require.False(t, rt.connCache.lru.Contains(inst.ID+"file"+fmt.Sprintf("dsn:%s ", repodsn)))
			_, ok := rt.migrationMetaCache.cache.Get(inst.ID)
			require.False(t, ok)
			_, err = svc.Olap.Execute(context.Background(), &drivers.Statement{Query: "SELECT COUNT(*) FROM rill.migration_version"})
			require.True(t, err != nil)

			// verify db file is dropped if requested
			_, err = os.Stat(dbFile)
			require.Equal(t, tt.dropDB, os.IsNotExist(err))
		})
	}
}

func TestRuntime_DeleteInstance_DropCorrupted(t *testing.T) {
	// We require the ability to delete instances and drop database files created with old versions of DuckDB, which can no longer be opened.

	// Prepare
	ctx := context.Background()
	rt := NewTestRunTime(t)
	dbpath := filepath.Join(t.TempDir(), "test.db")

	// Create instance
	inst := &drivers.Instance{
		ID:            "default",
		OLAPConnector: "duckdb",
		RepoConnector: "repo",
		EmbedCatalog:  true,
		Connectors: []*runtimev1.Connector{
			{
				Type:   "file",
				Name:   "repo",
				Config: map[string]string{"dsn": t.TempDir()},
			},
			{
				Type:   "duckdb",
				Name:   "duckdb",
				Config: map[string]string{"dsn": dbpath},
			},
		},
	}
	err := rt.CreateInstance(context.Background(), inst)
	require.NoError(t, err)

	// Put some data into it to create a .db file on disk
	olap, release, err := rt.OLAP(ctx, inst.ID)
	require.NoError(t, err)
	defer release()
	err = olap.Exec(ctx, &drivers.Statement{Query: "CREATE TABLE data(id INTEGER, name VARCHAR)"})
	require.NoError(t, err)

	// Close OLAP connection
	driver, cfg, err := rt.connectorConfig(ctx, inst.ID, inst.OLAPConnector)
	require.NoError(t, err)
	evicted := rt.connCache.evict(ctx, inst.ID, driver, cfg)
	require.True(t, evicted)

	// Corrupt database file
	err = os.WriteFile(dbpath, []byte("corrupted"), 0644)
	require.NoError(t, err)

	// Check we can't open it anymore
	_, _, err = rt.OLAP(ctx, inst.ID)
	require.Error(t, err)
	require.FileExists(t, dbpath)

	// Delete instance and check it still drops the .db file
	err = rt.DeleteInstance(ctx, inst.ID, true)
	require.NoError(t, err)
	require.NoFileExists(t, dbpath)
}

// New returns a runtime configured for use in tests.
func NewTestRunTime(t *testing.T) *Runtime {
	globalConnectors := []*runtimev1.Connector{
		{
			Type: "sqlite",
			Name: "metastore",
			// Setting a test-specific name ensures a unique connection when "cache=shared" is enabled.
			// "cache=shared" is needed to prevent threading problems.
			Config: map[string]string{"dsn": "file:rill?mode=memory&cache=shared"},
		},
	}

	opts := &Options{
		ConnectionCacheSize:     100,
		MetastoreConnector:      "metastore",
		QueryCacheSizeBytes:     int64(datasize.MB) * 100,
		SecurityEngineCacheSize: 100,
		AllowHostAccess:         true,
		SystemConnectors:        globalConnectors,
	}
	rt, err := New(opts, zap.NewNop(), activity.NewNoopClient())
	t.Cleanup(func() {
		rt.Close()
	})
	require.NoError(t, err)

	return rt
}

func connectorsEqual(a, b *runtimev1.Connector) bool {
	if (a != nil) != (b != nil) {
		return false
	}
	return a.Name == b.Name && a.Type == b.Type && maps.Equal(a.Config, b.Config)
}
