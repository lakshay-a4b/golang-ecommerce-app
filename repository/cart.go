// In repository/cart_repository.go
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Cart struct {
	CartID      int             `json:"cartId"`
	UserID      string          `json:"userId"`
	ProductInfo json.RawMessage `json:"productInfo"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type CartRepository struct {
	pool *pgxpool.Pool
}

func NewCartRepository(pool *pgxpool.Pool) *CartRepository {
	return &CartRepository{pool: pool}
}

func (r *CartRepository) GetCartByUserID(ctx context.Context, userID string) (*Cart, error) {
	query := `SELECT "cartId", "userId", "productInfo", "updatedAt" FROM cart WHERE "userId" = $1`
	row := r.pool.QueryRow(ctx, query, userID)

	var cart Cart
	err := row.Scan(&cart.CartID, &cart.UserID, &cart.ProductInfo, &cart.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error in GetCartByUserID (userId: %s): %v", userID, err)
		return nil, errors.New("failed to retrieve cart")
	}
	return &cart, nil
}

func (r *CartRepository) CreateOrUpdateCart(ctx context.Context, userID string, productInfo json.RawMessage) (*Cart, error) {
	query := `
		INSERT INTO cart ("userId", "productInfo")
		VALUES ($1, $2)
		ON CONFLICT ("userId")
		DO UPDATE SET 
			"productInfo" = EXCLUDED."productInfo",
			"updatedAt" = NOW()
		RETURNING "cartId", "userId", "productInfo", "updatedAt"
	`

	row := r.pool.QueryRow(ctx, query, userID, productInfo)

	var cart Cart
	err := row.Scan(&cart.CartID, &cart.UserID, &cart.ProductInfo, &cart.UpdatedAt)
	if err != nil {
		log.Printf("Error in CreateOrUpdateCart (userId: %s): %v", userID, err)
		return nil, errors.New("failed to create/update cart")
	}
	return &cart, nil
}

func (r *CartRepository) DeleteCart(ctx context.Context, userID string, tx pgx.Tx) error {
	query := `DELETE FROM cart WHERE "userId" = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		log.Printf("Error in DeleteCart (userId: %s): %v", userID, err)
		return errors.New("failed to delete cart")
	}
	return nil
}

func (r *CartRepository) UpdateCartProducts(ctx context.Context, userID string, productInfo json.RawMessage) (*Cart, error) {
	query := `
		UPDATE cart
		SET "productInfo" = $2,
			"updatedAt" = NOW()
		WHERE "userId" = $1
		RETURNING "cartId", "userId", "productInfo", "updatedAt"
	`

	row := r.pool.QueryRow(ctx, query, userID, productInfo)

	var cart Cart
	err := row.Scan(&cart.CartID, &cart.UserID, &cart.ProductInfo, &cart.UpdatedAt)
	if err != nil {
		log.Printf("Error in UpdateCartProducts (userId: %s): %v", userID, err)
		return nil, errors.New("failed to update cart products")
	}
	return &cart, nil
}