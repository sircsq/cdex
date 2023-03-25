package db

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Storage interface {
	CreateCollection(ctx context.Context, arg CreateCollectionParams) (Collection, error)
}

type NartDB struct {
	db *bun.DB
}

func NewNartDB(dsn string) *NartDB {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	return &NartDB{db: db}
}

func (db *NartDB) CreateCollection(ctx context.Context, arg CreateCollectionParams) (Collection, error) {
	return Collection{}, nil
}
