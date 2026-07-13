package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/eneipereira/go-order-service/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgCustomerRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCustomerRepository(pool *pgxpool.Pool) *pgCustomerRepository {
	return &pgCustomerRepository{pool: pool}
}

func (repo *pgCustomerRepository) Save(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	query := `INSERT INTO customers (name, email, phone, password_hash) VALUES ($1, $2, $3, $4) RETURNING id, name, email, phone, created_at, updated_at`

	err := repo.pool.QueryRow(
		ctx, query, customer.Name, customer.Email, customer.Phone, customer.PasswordHash,
	).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		return nil, HandleRepositoryError(err)
	}
	return customer, nil
}

func (repo *pgCustomerRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM customers ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := repo.pool.Query(ctx, query, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("Error listing customers: %w", err)
	}

	defer rows.Close()

	//customers, err := pgx.CollectRows(rows, pgx.RowToAddrOf[model.Customer])
	customers, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*model.Customer, error) {
        var c model.Customer
        
        err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.CreatedAt, &c.UpdatedAt)
        if err != nil {
            return nil, err
        }
        
        return &c, nil
    })

	if err != nil {
		return nil, fmt.Errorf("Error processing customer list: %w", err)
	}

	return customers, nil
}

func (repo *pgCustomerRepository) FindByID(ctx context.Context, id string) (*model.Customer, error) {
	query := `SELECT id, name, email, phone, created_at, updated_at FROM customers WHERE id = $1::uuid`

	customer := &model.Customer{}

	err := repo.pool.QueryRow(ctx, query, id).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("Error searching customer by id: %w", err)
	}
	return customer, nil
}

func HandleRepositoryError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return model.ErrEmailAlreadyExists
		case "23503":
			return model.ErrForeignKeyViolation
		case "23502":
			return fmt.Errorf("%w: %s", model.ErrNotNullViolation, pgErr.ColumnName)
		}
	}
	return fmt.Errorf("database error: %w", err)
}
