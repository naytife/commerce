// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createShop = `-- name: CreateShop :one
INSERT INTO shops (owner_id, title, default_domain, favicon_url,logo_url,email, currency_code, about, status, address,phone_number, seo_description, seo_keywords, seo_title)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING shop_id, owner_id, title, default_domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at
`

type CreateShopParams struct {
	OwnerID        uuid.UUID
	Title          string
	DefaultDomain  string
	FaviconUrl     pgtype.Text
	LogoUrl        pgtype.Text
	Email          string
	CurrencyCode   string
	About          pgtype.Text
	Status         string
	Address        pgtype.Text
	PhoneNumber    pgtype.Text
	SeoDescription pgtype.Text
	SeoKeywords    []string
	SeoTitle       pgtype.Text
}

func (q *Queries) CreateShop(ctx context.Context, arg CreateShopParams) (Shop, error) {
	row := q.db.QueryRow(ctx, createShop,
		arg.OwnerID,
		arg.Title,
		arg.DefaultDomain,
		arg.FaviconUrl,
		arg.LogoUrl,
		arg.Email,
		arg.CurrencyCode,
		arg.About,
		arg.Status,
		arg.Address,
		arg.PhoneNumber,
		arg.SeoDescription,
		arg.SeoKeywords,
		arg.SeoTitle,
	)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.DefaultDomain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const createShopCategory = `-- name: CreateShopCategory :one
INSERT INTO categories (slug, title, description, parent_id, shop_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING category_id, slug, title, description, parent_id, allowed_attributes, created_at, updated_at, shop_id
`

type CreateShopCategoryParams struct {
	Slug        string
	Title       string
	Description pgtype.Text
	ParentID    pgtype.Int8
	ShopID      int64
}

func (q *Queries) CreateShopCategory(ctx context.Context, arg CreateShopCategoryParams) (Category, error) {
	row := q.db.QueryRow(ctx, createShopCategory,
		arg.Slug,
		arg.Title,
		arg.Description,
		arg.ParentID,
		arg.ShopID,
	)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.ParentID,
		&i.AllowedAttributes,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ShopID,
	)
	return i, err
}

const getShop = `-- name: GetShop :one
SELECT shop_id, owner_id, title, default_domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at FROM shops
WHERE shop_id = $1
`

func (q *Queries) GetShop(ctx context.Context, shopID int64) (Shop, error) {
	row := q.db.QueryRow(ctx, getShop, shopID)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.DefaultDomain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getShopByDomain = `-- name: GetShopByDomain :one
SELECT shop_id, owner_id, title, default_domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at FROM shops
WHERE default_domain = $1
`

func (q *Queries) GetShopByDomain(ctx context.Context, defaultDomain string) (Shop, error) {
	row := q.db.QueryRow(ctx, getShopByDomain, defaultDomain)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.DefaultDomain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getShopCategory = `-- name: GetShopCategory :one
SELECT category_id, slug, title, description, parent_id, allowed_attributes, created_at, updated_at, shop_id FROM categories
WHERE category_id = $1
`

func (q *Queries) GetShopCategory(ctx context.Context, categoryID int64) (Category, error) {
	row := q.db.QueryRow(ctx, getShopCategory, categoryID)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.ParentID,
		&i.AllowedAttributes,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ShopID,
	)
	return i, err
}

const getShopsByOwner = `-- name: GetShopsByOwner :many
SELECT shop_id, owner_id, title, default_domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at FROM shops
WHERE owner_id = $1
`

func (q *Queries) GetShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]Shop, error) {
	rows, err := q.db.Query(ctx, getShopsByOwner, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shop
	for rows.Next() {
		var i Shop
		if err := rows.Scan(
			&i.ShopID,
			&i.OwnerID,
			&i.Title,
			&i.DefaultDomain,
			&i.FaviconUrl,
			&i.LogoUrl,
			&i.Email,
			&i.CurrencyCode,
			&i.Status,
			&i.About,
			&i.Address,
			&i.PhoneNumber,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.SeoTitle,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT user_id, auth0_sub, email, name, profile_picture_url, created_at, last_login FROM users
WHERE auth0_sub = $1
`

func (q *Queries) GetUser(ctx context.Context, auth0Sub pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUser, auth0Sub)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Auth0Sub,
		&i.Email,
		&i.Name,
		&i.ProfilePictureUrl,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}

const updateShop = `-- name: UpdateShop :one
UPDATE shops
SET 
    title = COALESCE($1, title),
    favicon_url = COALESCE($2, favicon_url),
    currency_code = COALESCE($3, currency_code),
    about = COALESCE($4, about),
    status = COALESCE($5, status),
    phone_number = COALESCE($6, phone_number),
    seo_description = COALESCE($7, seo_description),
    seo_keywords = COALESCE($8, seo_keywords),
    seo_title = COALESCE($9, seo_title),
    address = COALESCE($10, address),
    email = COALESCE($11, email)
WHERE default_domain = $12
RETURNING shop_id, owner_id, title, default_domain, favicon_url, logo_url, email, currency_code, status, about, address, phone_number, seo_description, seo_keywords, seo_title, updated_at, created_at
`

type UpdateShopParams struct {
	Title          pgtype.Text
	FaviconUrl     pgtype.Text
	CurrencyCode   pgtype.Text
	About          pgtype.Text
	Status         pgtype.Text
	PhoneNumber    pgtype.Text
	SeoDescription pgtype.Text
	SeoKeywords    []string
	SeoTitle       pgtype.Text
	Address        pgtype.Text
	Email          pgtype.Text
	DefaultDomain  string
}

func (q *Queries) UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error) {
	row := q.db.QueryRow(ctx, updateShop,
		arg.Title,
		arg.FaviconUrl,
		arg.CurrencyCode,
		arg.About,
		arg.Status,
		arg.PhoneNumber,
		arg.SeoDescription,
		arg.SeoKeywords,
		arg.SeoTitle,
		arg.Address,
		arg.Email,
		arg.DefaultDomain,
	)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.Title,
		&i.DefaultDomain,
		&i.FaviconUrl,
		&i.LogoUrl,
		&i.Email,
		&i.CurrencyCode,
		&i.Status,
		&i.About,
		&i.Address,
		&i.PhoneNumber,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.SeoTitle,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const upsertUser = `-- name: UpsertUser :one
INSERT INTO users (auth0_sub, email, name, profile_picture_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (auth0_sub)
DO UPDATE SET email = EXCLUDED.email, name = EXCLUDED.name, profile_picture_url = EXCLUDED.profile_picture_url
RETURNING user_id, auth0_sub, email, name, profile_picture_url
`

type UpsertUserParams struct {
	Auth0Sub          pgtype.Text
	Email             string
	Name              pgtype.Text
	ProfilePictureUrl pgtype.Text
}

type UpsertUserRow struct {
	UserID            uuid.UUID
	Auth0Sub          pgtype.Text
	Email             string
	Name              pgtype.Text
	ProfilePictureUrl pgtype.Text
}

func (q *Queries) UpsertUser(ctx context.Context, arg UpsertUserParams) (UpsertUserRow, error) {
	row := q.db.QueryRow(ctx, upsertUser,
		arg.Auth0Sub,
		arg.Email,
		arg.Name,
		arg.ProfilePictureUrl,
	)
	var i UpsertUserRow
	err := row.Scan(
		&i.UserID,
		&i.Auth0Sub,
		&i.Email,
		&i.Name,
		&i.ProfilePictureUrl,
	)
	return i, err
}
