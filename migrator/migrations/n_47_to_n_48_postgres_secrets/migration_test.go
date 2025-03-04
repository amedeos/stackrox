//go:build sql_integration

package n47ton48

import (
	"context"
	"strconv"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_47_to_n_48_postgres_secrets/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_47_to_n_48_postgres_secrets/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	ctx context.Context

	legacyDB   *rocksdb.RocksDB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	s.T().Setenv(env.PostgresDatastoreEnabled.EnvVar(), "true")
	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	var err error
	s.legacyDB, err = rocksdb.NewTemp(s.T().Name())
	s.NoError(err)

	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), true)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestSecretMigration() {
	newStore := pgStore.New(s.postgresDB.DB)
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)

	// Prepare data and write to legacy DB
	var secrets []*storage.Secret
	countBadIDs := 0
	for i := 0; i < 200; i++ {
		secret := &storage.Secret{}
		s.NoError(testutils.FullInit(secret, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		if i%10 == 0 {
			secret.Id = strconv.Itoa(i)
			countBadIDs = countBadIDs + 1
		}
		secrets = append(secrets, secret)
	}

	s.NoError(legacyStore.UpsertMany(s.ctx, secrets))

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(secrets)-countBadIDs, count)
	for _, secret := range secrets {
		if pgutils.NilOrUUID(secret.GetId()) != nil {
			fetched, exists, err := newStore.Get(s.ctx, secret.GetId())
			s.NoError(err)
			s.True(exists)
			s.Equal(secret, fetched)
		}
	}
}
