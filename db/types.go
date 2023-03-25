package db

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type Storage interface {
	Insert(ctx context.Context, value interface{}) error
	GetCollections(ctx context.Context, page, pageSize int) ([]Collection, error)
	GetCollectionByID(ctx context.Context, id int) (Collection, error)
	GetCollectionByCreator(ctx context.Context, address string, page, pageSize int) ([]Collection, error)
	GetCollectionItems(ctx context.Context, id, page, pageSize int) ([]Item, error)
}

type NartDB struct {
	db *bun.DB
}

func NewNartDB(dsn string) *NartDB {
	pgConn := pgdriver.NewConnector(pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(5*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second))

	sqlDB := sql.OpenDB(pgConn)
	db := bun.NewDB(sqlDB, pgdialect.New())

	return &NartDB{db: db}
}

func (db *NartDB) Insert(ctx context.Context, value interface{}) error {
	switch values := value.(type) {
	case []interface{}:
		for v := range values {
			if _, err := db.db.NewInsert().Model(v).Exec(ctx); err != nil {
				return err
			}
		}
	case interface{}:
		if _, err := db.db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (db *NartDB) GetCollections(ctx context.Context, page, pageSize int) ([]Collection, error) {
	var collections []Collection
	err := db.db.NewSelect().Model(&collections).Limit(pageSize).Offset(pageSize * (page - 1)).Scan(ctx)
	return collections, err
}

func (db *NartDB) GetCollectionByCreator(ctx context.Context, address string, page, pageSize int) ([]Collection, error) {
	var collections []Collection
	err := db.db.NewSelect().Model(&collections).Where("creator = ?", address).Limit(pageSize).Offset(pageSize * (page - 1)).Scan(ctx)
	return collections, err
}

func (db *NartDB) GetCollectionByID(ctx context.Context, id int) (Collection, error) {
	var c Collection
	err := db.db.NewSelect().Model(&c).Where("id = ?", id).Scan(ctx)
	return c, err
}

func (db *NartDB) GetCollectionItems(ctx context.Context, id, page, pageSize int) ([]Item, error) {
	var items []Item
	err := db.db.NewSelect().Model(&items).Where("collection_id = ?", id).Limit(pageSize).Offset(pageSize * (page - 1)).Scan(ctx)
	return items, err
}
