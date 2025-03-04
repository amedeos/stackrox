package n31ton32

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	frozenSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_31_to_n_32_postgres_network_graph_configs/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_31_to_n_32_postgres_network_graph_configs/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/sac"
	"gorm.io/gorm"
)

const (
	networkGraphConfigKey = "networkGraphConfig"
)

var (
	startingSeqNum = pkgMigrations.BasePostgresDBVersionSeqNum() + 31 // 142

	migration = types.Migration{
		StartingSeqNum: startingSeqNum,
		VersionAfter:   &storage.Version{SeqNum: int32(startingSeqNum + 1)}, // 143
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving network_graph_configs from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = frozenSchema.NetworkGraphConfigsSchema
	log       = loghelper.LogWrapper{}
)

func move(gormDB *gorm.DB, postgresDB *postgres.DB, legacyStore legacy.Store) error {
	ctx := sac.WithAllAccess(context.Background())
	store := pgStore.New(postgresDB)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableNetworkGraphConfigsStmt)
	var networkGraphConfigs []*storage.NetworkGraphConfig

	var found bool
	err := walk(ctx, legacyStore, func(obj *storage.NetworkGraphConfig) error {
		if found {
			log.WriteToStderr("found multiple network graph configs")
			return nil
		}
		found = true
		obj.Id = networkGraphConfigKey
		networkGraphConfigs = append(networkGraphConfigs, obj)
		if len(networkGraphConfigs) == batchSize {
			if err := store.UpsertMany(ctx, networkGraphConfigs); err != nil {
				log.WriteToStderrf("failed to persist network_graph_configs to store %v", err)
				return err
			}
			networkGraphConfigs = networkGraphConfigs[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(networkGraphConfigs) > 0 {
		if err = store.UpsertMany(ctx, networkGraphConfigs); err != nil {
			log.WriteToStderrf("failed to persist network_graph_configs to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.NetworkGraphConfig) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
