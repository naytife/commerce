package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repoSvc struct {
	*Queries
	db *pgxpool.Pool
}

type contextKey string

const shopIDKey contextKey = "shop_id"

func getShopIDFromContext(ctx context.Context) (int, error) {
	// Retrieve the shopID from the context
	shopID, ok := ctx.Value(shopIDKey).(int)
	if !ok {
		return 0, errors.New("shop_id not found in context")
	}
	return shopID, nil
}
func setShopIDInSession(ctx context.Context, r *repoSvc, shopID int) error {
	// Your logic to set the shop_id, e.g., using a SQL command
	_, err := r.db.Exec(ctx, "SET LOCAL shop_id = $1", shopID)
	return err
}
func (r *repoSvc) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	// Set shop_id for the session
	log.Println("HELO HTER")
	shopID, err := getShopIDFromContext(ctx) // Assume this function returns an int
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to retrieve shop_id from context: %w", err)
	}

	if err := setShopIDInSession(ctx, r, shopID); err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to set shop_id in session: %w", err)
	}

	// Execute the original sqlc query
	commandTag, err := r.db.Exec(ctx, sql, arguments...)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("execution failed: %w", err)
	}
	return commandTag, nil
}

// wrap the Query method to set shop_id before executing any query
func (r *repoSvc) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	// Set shop_id for the session
	log.Println("Executing SQL:", sql)
	shopID, err := getShopIDFromContext(ctx) // Assume this function returns an int
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shop_id from context: %w", err)
	}

	if err := setShopIDInSession(ctx, r, shopID); err != nil {
		return nil, fmt.Errorf("failed to set shop_id in session: %w", err)
	}

	// Execute the original sqlc query
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return rows, nil
}
func (r *repoSvc) withTx(ctx context.Context, txFn func(*Queries) error) error {
	log.Println("IN WITHTX")
	shopID, ok := ctx.Value(shopIDKey).(int)

	if !ok {
		return fmt.Errorf("shop_id not found in context")
	}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "SET LOCAL shop_id = $1", shopID)
	if err != nil {
		return fmt.Errorf("failed to set shop_id: %w", err)
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
	// USER
	UpsertUser(ctx context.Context, arg UpsertUserParams) (UpsertUserRow, error)
	GetUser(ctx context.Context, auth0Sub pgtype.Text) (User, error)
	// SHOP
	CreateShop(ctx context.Context, shopArg CreateShopParams) (Shop, error)
	UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error)
	GetShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]Shop, error)
	GetShopByDomain(ctx context.Context, defaultDomain string) (Shop, error)
	GetShopIDByDomain(ctx context.Context, domain string) (int64, error)
	// CATEGORY
	CreateShopCategory(ctx context.Context, arg CreateShopCategoryParams) (Category, error)
	GetShopCategory(ctx context.Context, categoryID int64) (Category, error)
	UpdateShopCategory(ctx context.Context, arg UpdateShopCategoryParams) (Category, error)
	GetShopCategories(ctx context.Context) ([]Category, error)
	CreateCategoryAttribute(ctx context.Context, arg CreateCategoryAttributeParams) ([]byte, error)
	DeleteCategoryAttribute(ctx context.Context, arg DeleteCategoryAttributeParams) ([]byte, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

func InitDB(dataSourceName string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dataSourceName)
	config.MaxConns = 50
	config.MinConns = 5
	config.MaxConnIdleTime = 5 * time.Minute
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
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	shop := Shop{}
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		shop, err = q.CreateShop(ctx, shopArg)
		return err
	})
	return shop, err
}

func (r *repoSvc) UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	shop := Shop{}
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		shop, err = q.UpdateShop(ctx, arg)
		if err != nil {
			log.Println(err)
		}
		return err
	})
	return shop, err
}

func (r *repoSvc) CreateShopCategory(ctx context.Context, arg CreateShopCategoryParams) (Category, error) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	category := Category{}
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		category, err = q.CreateShopCategory(ctx, arg)
		return err
	})
	return category, err
}
