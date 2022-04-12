// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"reflect"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/metrics"
	pkgSchema "github.com/stackrox/rox/central/postgres/schema"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

const (
	baseTable  = "image_cve_relations"
	countStmt  = "SELECT COUNT(*) FROM image_cve_relations"
	existsStmt = "SELECT EXISTS(SELECT 1 FROM image_cve_relations WHERE Id = $1 AND ImageId = $2 AND ImageCveId = $3)"

	getStmt    = "SELECT serialized FROM image_cve_relations WHERE Id = $1 AND ImageId = $2 AND ImageCveId = $3"
	deleteStmt = "DELETE FROM image_cve_relations WHERE Id = $1 AND ImageId = $2 AND ImageCveId = $3"
	walkStmt   = "SELECT serialized FROM image_cve_relations"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log    = logging.LoggerForModule()
	schema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable(baseTable)
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ImageCVEEdge)(nil)), baseTable).
			WithReference(func() *walker.Schema {
				parent := globaldb.GetSchemaForTable("images")
				if parent != nil {
					return parent
				}
				parent = walker.Walk(reflect.TypeOf((*storage.Image)(nil)), "images")
				globaldb.RegisterTable(parent)
				return parent
			}())
		globaldb.RegisterTable(schema)
		return schema
	}()
)

type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string, imageId string, imageCveId string) (bool, error)
	Get(ctx context.Context, id string, imageId string, imageCveId string) (*storage.ImageCVEEdge, bool, error)

	Walk(ctx context.Context, fn func(obj *storage.ImageCVEEdge) error) error

	AckKeysIndexed(ctx context.Context, keys ...string) error
	GetKeysToIndex(ctx context.Context) ([]string, error)
}

type storeImpl struct {
	db *pgxpool.Pool
}

// New returns a new Store instance using the provided sql instance.
func New(ctx context.Context, db *pgxpool.Pool) Store {
	pgutils.CreateTable(ctx, db, pkgSchema.CreateTableImagesStmt)
	pgutils.CreateTable(ctx, db, pkgSchema.CreateTableImageCveRelationsStmt)

	return &storeImpl{
		db: db,
	}
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Count, "ImageCVEEdge")

	row := s.db.QueryRow(ctx, countStmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string, imageId string, imageCveId string) (bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Exists, "ImageCVEEdge")

	row := s.db.QueryRow(ctx, existsStmt, id, imageId, imageCveId)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, pgutils.ErrNilIfNoRows(err)
	}
	return exists, nil
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string, imageId string, imageCveId string) (*storage.ImageCVEEdge, bool, error) {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Get, "ImageCVEEdge")

	conn, release, err := s.acquireConn(ctx, ops.Get, "ImageCVEEdge")
	if err != nil {
		return nil, false, err
	}
	defer release()

	row := conn.QueryRow(ctx, getStmt, id, imageId, imageCveId)
	var data []byte
	if err := row.Scan(&data); err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	var msg storage.ImageCVEEdge
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func (s *storeImpl) acquireConn(ctx context.Context, op ops.Op, typ string) (*pgxpool.Conn, func(), error) {
	defer metrics.SetAcquireDBConnDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.Release, nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.ImageCVEEdge) error) error {
	rows, err := s.db.Query(ctx, walkStmt)
	if err != nil {
		return pgutils.ErrNilIfNoRows(err)
	}
	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return err
		}
		var msg storage.ImageCVEEdge
		if err := proto.Unmarshal(data, &msg); err != nil {
			return err
		}
		if err := fn(&msg); err != nil {
			return err
		}
	}
	return nil
}

//// Used for testing

func dropTableImageCveRelations(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS image_cve_relations CASCADE")

}

func Destroy(ctx context.Context, db *pgxpool.Pool) {
	dropTableImageCveRelations(ctx, db)
}

//// Stubs for satisfying legacy interfaces

// AckKeysIndexed acknowledges the passed keys were indexed
func (s *storeImpl) AckKeysIndexed(ctx context.Context, keys ...string) error {
	return nil
}

// GetKeysToIndex returns the keys that need to be indexed
func (s *storeImpl) GetKeysToIndex(ctx context.Context) ([]string, error) {
	return nil, nil
}
