package db

import (
	"cdex/exchange"
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type Storage interface {
	Insert(ctx context.Context, value interface{}) error
	InsertCollection(ctx context.Context, arg CreateCollectionParams) (*Collection, error)
	InsertItem(ctx context.Context, arg CreateItemParams) (*Item, error)
	GetCollections(ctx context.Context, page, pageSize int) ([]*Collection, error)
	GetItems(ctx context.Context, page, pageSize int) ([]*Item, error)
	GetCollectionByID(ctx context.Context, id int) (*Collection, error)
	GetCollectionByCreator(ctx context.Context, address string, page, pageSize int) ([]*Collection, error)
	GetCollectionItems(ctx context.Context, id, page, pageSize int) ([]*Item, error)

	GetOrders(ctx context.Context, bid, page, pageSize int, status, sort string) ([]*exchange.Order, error)
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

func (db *NartDB) Insert(ctx context.Context, value interface{}) error {
	switch values := value.(type) {
	case []interface{}:
		for v := range values {
			if _, err := db.db.NewInsert().
				Model(v).
				Exec(ctx); err != nil {
				return err
			}
		}
	case interface{}:
		_, err := db.db.NewInsert().
			Model(value).
			Exec(ctx)
		return err
	}
	return nil
}

func (db *NartDB) InsertCollection(ctx context.Context, arg CreateCollectionParams) (*Collection, error) {
	var (
		err error
		c   Collection
	)

	err = db.db.QueryRowContext(ctx, "INSERT INTO collections(name,address,creator,chain,visible,status,created_at,type,tax,symbol,currency,image,background,banner,properties,introduction,description,twitter,instagram,discord,web) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) RETURNING *",
		&arg.Name, &arg.Address, &arg.Creator, &arg.Chain, &arg.Visible, &arg.Status, &arg.CreatedAt, &arg.Type, &arg.Tax, &arg.Symbol, &arg.Currency, &arg.Image, &arg.Background, &arg.Banner, &arg.Properties, &arg.Introduction, &arg.Description, &arg.Twitter, &arg.Instagram, &arg.Discord, &arg.Web).
		Scan(&c.ID, &c.Name, &c.Address, &c.Creator, &c.Chain, &c.Visible, &c.Status, &c.CreatedAt, &c.Type, &c.Tax, &c.Symbol, &c.Currency,
			&c.Image, &c.Background, &c.Banner, &c.Properties, &c.Introduction, &c.Description, &c.Twitter, &c.Instagram, &c.Discord, &c.Web)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *NartDB) InsertItem(ctx context.Context, arg CreateItemParams) (*Item, error) {
	var (
		err  error
		item Item
	)

	err = db.db.QueryRowContext(ctx, "INSERT INTO items(name,collection,token_id,creator,created_at,chain,image,description,properties) VALUES(?,?,?,?,?,?,?,?,?) RETURNING *",
		&arg.Name, &arg.Collection, &arg.TokenID, &arg.Creator, &arg.CreatedAt, &arg.Chain, &arg.Image, &arg.Description, &arg.Properties).
		Scan(&item.Name, &item.Collection, &item.TokenID, &item.Creator, &item.CreatedAt, &item.Chain, &item.Image, &item.Description, &item.Properties)
	if err != nil {
		return nil, err
	}

	return &item, nil
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
	err = db.db.NewSelect().Model((*Collection)(nil)).Order("id DESC").Limit(pageSize).Offset(pageSize*(page-1)).Scan(ctx, &collections)
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (db *NartDB) GetItems(ctx context.Context, page, pageSize int) ([]*Item, error) {
	var (
		err   error
		items []*Item
	)
	err = db.db.NewSelect().Model((*Item)(nil)).Order("created_at DESC").Limit(pageSize).Offset(pageSize*(page-1)).Scan(ctx, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *NartDB) GetCollectionItems(ctx context.Context, id, page, pageSize int) ([]*Item, error) {
	var items []*Item
	err := db.db.NewSelect().Model(&items).Where("collection = ?", id).Order("created_at DESC").Limit(pageSize).Offset(pageSize * (page - 1)).Scan(ctx)
	return items, err
}

func (db *NartDB) GetCollectionByCreator(ctx context.Context, address string, page, pageSize int) ([]*Collection, error) {
	var (
		err         error
		collections []*Collection
	)

	err = db.db.NewSelect().Model((*Collection)(nil)).Where("creator = ?", address).Order("id DESC").Limit(pageSize).Offset(pageSize*(page-1)).Scan(ctx, &collections)
	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (db *NartDB) GetOrders(ctx context.Context, bid, page, pageSize int, status, sort string) ([]*exchange.Order, error) {
	var (
		err    error
		orders []*exchange.Order
	)
	if sort == "desc" {
		err = db.db.NewSelect().Model((*exchange.Order)(nil)).Where("bid = ? AND status = ?", bid, status).Order("created_at DESC").Limit(pageSize).Offset(pageSize*(page-1)).Scan(ctx, &orders)
	} else {
		err = db.db.NewSelect().Model((*exchange.Order)(nil)).Where("bid = ? AND status = ?", bid, status).Order("created_at ASC").Limit(pageSize).Offset(pageSize*(page-1)).Scan(ctx, &orders)
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}
