package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
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
	defer tx.Rollback(ctx)

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
	SetShopIDInSession(ctx context.Context, shopID int64) error
	// USER
	UpsertUser(ctx context.Context, arg UpsertUserParams) (UpsertUserRow, error)
	GetUser(ctx context.Context, auth0Sub *string) (User, error)
	// SHOP
	CreateShop(ctx context.Context, shopArg CreateShopParams) (Shop, error)
	GetShop(ctx context.Context, shopID int64) (Shop, error)
	UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error)
	GetShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]Shop, error)
	GetShopByDomain(ctx context.Context, defaultDomain string) (Shop, error)
	GetShopIDByDomain(ctx context.Context, domain string) (int64, error)
	GetShopImages(ctx context.Context, shopID int64) (ShopImage, error)
	GetShopWhatsApp(ctx context.Context, shopID int64) (Whatsapp, error)
	GetShopFacebook(ctx context.Context, shopID int64) (Facebook, error)
	UpsertShopWhatsapp(ctx context.Context, arg UpsertShopWhatsappParams) (Whatsapp, error)
	UpsertShopFacebook(ctx context.Context, arg UpsertShopFacebookParams) (Facebook, error)
	// CATEGORY
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	GetCategory(ctx context.Context, arg GetCategoryParams) (GetCategoryRow, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	GetCategoryChildren(ctx context.Context, arg GetCategoryChildrenParams) ([]GetCategoryChildrenRow, error)
	GetCategories(ctx context.Context, arg GetCategoriesParams) ([]GetCategoriesRow, error)
	CreateCategoryAttribute(ctx context.Context, arg CreateCategoryAttributeParams) ([]byte, error)
	DeleteCategoryAttribute(ctx context.Context, arg DeleteCategoryAttributeParams) ([]byte, error)
	GetCategoryAttributes(ctx context.Context, categoryID int64) ([]byte, error)
	// PRODUCT
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	GetProducts(ctx context.Context, arg GetProductsParams) ([]GetProductsRow, error)
	GetProduct(ctx context.Context, arg GetProductParams) (GetProductRow, error)
	GetProductsByCategory(ctx context.Context, arg GetProductsByCategoryParams) ([]GetProductsByCategoryRow, error)
	GetProductAllowedAttributes(ctx context.Context, productID int64) ([]byte, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	CreateProductAllowedAttribute(ctx context.Context, arg CreateProductAllowedAttributeParams) ([]byte, error)
	DeleteProductAllowedAttribute(ctx context.Context, arg DeleteProductAllowedAttributeParams) ([]byte, error)
	UpsertProductVariations(ctx context.Context, shopID int64, productID int64, variations []UpsertProductVariationParams) ([]ProductVariation, error)
	GetProductVariations(ctx context.Context, arg GetProductVariationsParams) ([]ProductVariation, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repoSvc{
		Queries: New(db),
		db:      db,
	}
}

func InitDB(dataSourceName string) (*pgxpool.Pool, error) {
	// Initialize a new logger (using Logrus)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) // Set log level to Debug for detailed SQL logging

	// Wrap the Logrus logger in a TraceLog object
	traceLogger := &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
			switch level {
			case tracelog.LogLevelError:
				logger.WithFields(logrus.Fields(data)).Error(msg)
			case tracelog.LogLevelWarn:
				logger.WithFields(logrus.Fields(data)).Warn(msg)
			case tracelog.LogLevelInfo:
				logger.WithFields(logrus.Fields(data)).Info(msg)
			case tracelog.LogLevelDebug:
				logger.WithFields(logrus.Fields(data)).Debug(msg)
			}
		}),
		LogLevel: tracelog.LogLevelDebug, // Set log level to Debug
	}

	// Parse the pool config from the data source name
	config, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}

	// Attach the tracer to the config
	config.ConnConfig.Tracer = traceLogger

	// Set pool settings
	config.MaxConns = 1
	config.MinConns = 1
	config.MaxConnIdleTime = 5 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (r *repoSvc) SetShopIDInSession(ctx context.Context, shopID int64) error {
	// Assuming you want to execute a SQL command to set shop_id for the current session
	query := fmt.Sprintf("SET commerce.current_shop_id = %d", shopID)
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to set shop_id in session: %w", err)
	}
	return nil
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
		return err
	})
	return shop, err
}

func (r *repoSvc) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	category := Category{}
	err := r.withTx(ctx, func(q *Queries) error {
		var err error
		category, err = q.CreateCategory(ctx, arg)
		return err
	})
	return category, err
}

func (r *repoSvc) UpsertProductVariations(ctx context.Context, shopID int64, productID int64, variations []UpsertProductVariationParams) ([]ProductVariation, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	objsDB := []ProductVariation{}
	objsID := []int64{}

	err := r.withTx(ctx, func(q *Queries) error {
		// Batch upsert product variations
		batch := q.UpsertProductVariation(ctx, variations)

		// Process upserts
		batch.Query(func(i int, result []ProductVariation, err error) {
			if err != nil {
				fmt.Errorf("failed to upsert product variation: %w", err)
				return
			}

			for _, objDB := range result {
				objsID = append(objsID, objDB.ProductVariationID)
				objsDB = append(objsDB, objDB)
			}
		})

		if err := batch.Close(); err != nil {
			return fmt.Errorf("batch execution failed: %w", err)
		}

		// Batch delete the old variations with batchexec
		deleteBatch := q.DeleteProductVariations(ctx, []DeleteProductVariationsParams{
			{ShopID: shopID, ProductID: productID, ProductVariationIds: objsID},
		})

		// No need to process results here since :batchexec doesn’t return rows
		if err := deleteBatch.Close(); err != nil {
			return fmt.Errorf("batch delete failed: %w", err)
		}

		return nil
	})

	return objsDB, err
}
