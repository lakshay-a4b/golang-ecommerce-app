package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/models"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *OrderRepository) AddOrder(ctx context.Context, order models.Order, tx pgx.Tx) (*models.Order, error) {
	query := `
		INSERT INTO orders ("paymentId", "userId", "productInfo", status)
		VALUES ($1, $2, $3, $4)
		RETURNING "orderId", "paymentId", "userId", "productInfo", status, "createdAt"
	`

	var db pgx.Tx
	if tx != nil {
		db = tx
	} else {
		var err error
		db, err = r.pool.Begin(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer db.Rollback(ctx)
	}

	var newOrder models.Order
	err := db.QueryRow(ctx, query,
		order.PaymentID,
		order.UserId,
		order.ProductInfo,
		order.Status,
	).Scan(
		&newOrder.OrderID,
		&newOrder.PaymentID,
		&newOrder.UserId,
		&newOrder.ProductInfo,
		&newOrder.Status,
		&newOrder.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	if tx == nil {
		if err := db.Commit(ctx); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return &newOrder, nil
}

func (r *OrderRepository) GetOrdersByUser(ctx context.Context, userId string) ([]models.Order, error) {
	query := `
		SELECT "orderId", "paymentId", "userId", "productInfo", status, "createdAt"
		FROM orders
		WHERE "userId" = $1
		ORDER BY "createdAt" DESC
	`

	rows, err := r.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.OrderID,
			&o.PaymentID,
			&o.UserId,
			&o.ProductInfo,
			&o.Status,
			&o.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderId string) (*models.Order, error) {
	if orderId == "" {
		return nil, fmt.Errorf("invalid order ID")
	}

	query := `
		SELECT "orderId", "paymentId", "userId", "productInfo", status, "createdAt"
		FROM orders
		WHERE "orderId" = $1
	`

	var order models.Order
	err := r.pool.QueryRow(ctx, query, orderId).Scan(
		&order.OrderID,
		&order.PaymentID,
		&order.UserId,
		&order.ProductInfo,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}
func (r *OrderRepository) UpdateOrder(ctx context.Context, order models.Order) (*models.Order, error) {
	query := `
		UPDATE orders
		SET status = $1
		WHERE "orderId" = $2 AND "userId" = $3
		RETURNING "
		orderId", "paymentId", "userId", "productInfo", status, "createdAt"
	`
	var updatedOrder models.Order
	err := r.pool.QueryRow(ctx, query,
		order.Status,
		order.OrderID,
		order.UserId,
	).Scan(
		&updatedOrder.OrderID,
		&updatedOrder.PaymentID,
		&updatedOrder.UserId,
		&updatedOrder.ProductInfo,
		&updatedOrder.Status,
		&updatedOrder.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found or does not belong to user")
		}
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return &updatedOrder, nil
}