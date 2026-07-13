package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/eneipereira/go-order-service/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgProductRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresProductRepository(pool *pgxpool.Pool) *pgProductRepository {
	return &pgProductRepository{pool: pool}
}

func (repo *pgProductRepository) Save(ctx context.Context, product *model.Product) (*model.Product, error) {
	query := `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id, name, price, stock,created_at, updated_at`
	err := repo.pool.QueryRow(ctx, query, product.Name, product.Price, product.Stock).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, HandleRepositoryError(err)
	}
	return product, nil
}

func (repo *pgProductRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	query := `SELECT id, name, price, stock, created_at, updated_at FROM products ORDER BY name ASC LIMIT $1 OFFSET $2`
	rows, err := repo.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error listing products: %w", err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*model.Product, error) {
		var p model.Product
		err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &p, nil
	})

	return products, nil
}

func (repo *pgProductRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	query := `SELECT id, name, price, stock, created_at, updated_at FROM products WHERE id = $1::uuid`

	product := &model.Product{}

	err := repo.pool.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrProductNotFound
		}
		return nil, fmt.Errorf("Error searching product by id: %w", err)
	}
	return product, nil
}
