package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repoSvc struct {
	*Queries
	db *pgxpool.Pool
}

func (r *repoSvc) withTx(ctx context.Context, txFn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	q := New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx failed: %w, unable to rollback: %w", err, rbErr)
		}
	} else {
		err = tx.Commit(ctx)
	}
	return err
}

type Repository interface {
	CreateShop(ctx context.Context, shopArg CreateShopParams) (Shop, error)
	GetShop(ctx context.Context, id pgtype.UUID) (Shop, error)
	GetShopsByOwner(ctx context.Context, ownerID pgtype.UUID) ([]Shop, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

func Open(dataSourceName string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.New(context.Background(), config.ConnString())
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (r *repoSvc) CreateShop(ctx context.Context, shopArg CreateShopParams) (Shop, error) {
	shop := Shop{}
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		shop, err = q.CreateShop(ctx, shopArg)
		return err
	})
	return shop, err
}
