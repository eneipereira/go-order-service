package repository

import (
	"context"
	"fmt"

	"github.com/eneipereira/go-order-service/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) *pgOrderRepository {
	return &pgOrderRepository{pool: pool}
}

func (r *pgOrderRepository) Create(ctx context.Context, order *model.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Error beginning transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	orderQuery := `INSERT INTO orders (customer_id, status, total) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err = tx.QueryRow(ctx, orderQuery, order.CustomerID, order.Status, order.Total).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Error inserting order: %w", HandleRepositoryError(err))
	}

	for i := range order.Items {
		item := &order.Items[i]
		item.OrderID = order.ID

		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`
		err = tx.QueryRow(ctx, itemQuery, item.OrderID, item.ProductID, item.Quantity, item.Price).Scan(&item.ID)
		if err != nil {
			return fmt.Errorf("Error inserting order item: %w", HandleRepositoryError(err))
		}

		stockUpdateQuery := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
		res, err := tx.Exec(ctx, stockUpdateQuery, item.Quantity, item.ProductID)
		if err != nil {
			return fmt.Errorf("Error updating stock: %w", err)
		}
		if res.RowsAffected() == 0 {
			return model.ErrInsufficientStock
		}
	}

	return tx.Commit(ctx)
}

func (r *pgOrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {

	orderQuery := `SELECT id, customer_id, status, total, created_at, updated_at FROM orders WHERE id = $1::uuid`
	order := &model.Order{}
	err := r.pool.QueryRow(ctx, orderQuery, id).Scan(&order.ID, &order.CustomerID, &order.Status, &order.Total, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, HandleRepositoryError(err)
	}

	itemsQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1::uuid`
	rows, err := r.pool.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("Error fetching order items: %w", err)
	}
	defer rows.Close()

	order.Items, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.OrderItem])
	if err != nil {
		return nil, fmt.Errorf("Error collecting order items: %w", err)
	}

	return order, nil
}

func (r *pgOrderRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Order, error) {

	ordersQuery := `SELECT id, customer_id, status, total, created_at, updated_at FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, ordersQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error listing all orders: %w", err)
	}
	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToAddrOf[model.Order])
	if err != nil {
		return nil, fmt.Errorf("Error collecting orders: %w", err)
	}

	if len(orders) == 0 {
		return []*model.Order{}, nil
	}

	orderIDs := make([]uuid.UUID, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	itemsQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = ANY($1::uuid[])`
	itemRows, err := r.pool.Query(ctx, itemsQuery, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("Error fetching items for orders: %w", err)
	}
	defer itemRows.Close()

	items, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[model.OrderItem])
	if err != nil {
		return nil, fmt.Errorf("Error collecting items for orders: %w", err)
	}

	itemsByOrderID := make(map[uuid.UUID][]model.OrderItem)
	for _, item := range items {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}
	for _, order := range orders {
		order.Items = itemsByOrderID[order.ID]
	}

	return orders, nil
}

func (r *pgOrderRepository) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2::uuid`
	res, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("Error updating order status: %w", err)
	}
	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}
	return nil
}

func (r *pgOrderRepository) CancelOrderAndRestockItems(ctx context.Context, order *model.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Error beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	orderUpdateQuery := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2::uuid`
	_, err = tx.Exec(ctx, orderUpdateQuery, model.StatusCanceled, order.ID)
	if err != nil {
		return fmt.Errorf("Error canceling order: %w", err)
	}

	for _, item := range order.Items {
		stockUpdateQuery := `UPDATE products SET stock = stock + $1 WHERE id = $2::uuid`
		_, err := tx.Exec(ctx, stockUpdateQuery, item.Quantity, item.ProductID)
		if err != nil {
			return fmt.Errorf("Error restoking product %s: %w", item.ProductID, err)
		}
	}

	return tx.Commit(ctx)
}
