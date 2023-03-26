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
	InsertCollection(ctx context.Context, arg CreateCollectionParams) (*Collection, error)
	GetCollections(ctx context.Context, page, pageSize int) ([]*Collection, error)
	GetCollectionByID(ctx context.Context, id int) (*Collection, error)
	GetCollectionByCreator(ctx context.Context, address string, page, pageSize int) ([]*Collection, error)
	GetCollectionItems(ctx context.Context, id, page, pageSize int) ([]*Item, error)
}

type NartDB struct {
	db *bun.DB
}

func NewNartDB(dsn string) *NartDB {
	pgConn := pgdriver.NewConnector(pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(10*time.Second),
		pgdriver.WithDialTimeout(10*time.Second),
		pgdriver.WithReadTimeout(10*time.Second),
		pgdriver.WithWriteTimeout(10*time.Second))

	sqlDB := sql.OpenDB(pgConn)
	db := bun.NewDB(sqlDB, pgdialect.New())

	return &NartDB{db: db}
}

func (db *NartDB) InsertCollection(ctx context.Context, arg CreateCollectionParams) (*Collection, error) {
	var (
		err error
		c   Collection
	)

	err = db.db.QueryRowContext(ctx, "insert into collections(name,address,creator,chain,visible,status,created_at,type,tax,symbol,currency,image,background,banner,properties,introduction,description,twitter,instagram,discord,web) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) returning *",
		&arg.Name, &arg.Address, &arg.Creator, &arg.Chain, &arg.Visible, &arg.Status, &arg.CreatedAt, &arg.Type, &arg.Tax, &arg.Symbol, &arg.Currency, &arg.Image, &arg.Background, &arg.Banner, &arg.Properties, &arg.Introduction, &arg.Description, &arg.Twitter, &arg.Instagram, &arg.Discord, &arg.Web).
		Scan(&c.ID, &c.Name, &c.Address, &c.Creator, &c.Chain, &c.Visible, &c.Status, &c.CreatedAt, &c.Type, &c.Tax, &c.Symbol, &c.Currency,
			&c.Image, &c.Background, &c.Banner, &c.Properties, &c.Introduction, &c.Description, &c.Twitter, &c.Instagram, &c.Discord, &c.Web)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *NartDB) GetCollectionByID(ctx context.Context, id int) (*Collection, error) {
	var c Collection
	err := db.db.NewSelect().Model(&c).Where("id = ?", id).Scan(ctx)
	return &c, err
}

func (db *NartDB) GetCollections(ctx context.Context, page, pageSize int) ([]*Collection, error) {
	var (
		err         error
		collections []*Collection
	)
	err = db.db.NewSelect().Model((*Collection)(nil)).Limit(pageSize).Offset(pageSize*(page-1)).Scan(ctx, &collections)
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (db *NartDB) GetCollectionItems(ctx context.Context, id, page, pageSize int) ([]*Item, error) {
	var items []*Item
	err := db.db.NewSelect().Model(&items).Where("collection_id = ?", id).Limit(pageSize).Offset(pageSize * (page - 1)).Scan(ctx)
	return items, err
}

func (db *NartDB) GetCollectionByCreator(ctx context.Context, address string, page, pageSize int) ([]*Collection, error) {
	return nil, nil
}
