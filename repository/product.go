package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/models"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// GetAllProducts fetches all products
func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT "productId", name, description, image, price, "createdAt" FROM products`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Printf("Database error: GetAllProducts failed: %v", err)
		return nil, fmt.Errorf("failed to fetch all products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ProductID,
			&p.Name,
			&p.Description,
			&p.Image,
			&p.Price,
			&p.CreatedAt,
		); err != nil {
			log.Printf("Row scan error in GetAllProducts: %v", err)
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error in GetAllProducts: %w", err)
	}

	return products, nil
}

// GetPaginatedProducts fetches products by limit and offset
func (r *ProductRepository) GetPaginatedProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {
	query := `SELECT "productId", name, description, image, price, "createdAt" 
	          FROM products ORDER BY "productId" LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		log.Printf("Database error: GetPaginatedProducts failed: %v", err)
		return nil, fmt.Errorf("failed to fetch paginated products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ProductID,
			&p.Name,
			&p.Description,
			&p.Image,
			&p.Price,
			&p.CreatedAt,
		); err != nil {
			log.Printf("Row scan error in GetPaginatedProducts: %v", err)
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error in GetPaginatedProducts: %w", err)
	}

	return products, nil
}

// GetTotalProductCount returns total count of products
func (r *ProductRepository) GetTotalProductCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM products`

	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("Database error: GetTotalProductCount failed: %v", err)
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// GetProductByID fetches a single product by its ID
func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	query := `SELECT "productId", name, description, image, price, "createdAt"
	          FROM products WHERE "productId" = $1`

	var p models.Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ProductID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Database error: GetProductByID(%d) failed: %v", id, err)
		return nil, fmt.Errorf("failed to fetch product by ID: %w", err)
	}

	return &p, nil
}

// CreateProduct inserts a new product
func (r *ProductRepository) CreateProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	query := `
		INSERT INTO products (name, description, image, price)
		VALUES ($1, $2, $3, $4)
		RETURNING "productId", name, description, image, price, "createdAt"
	`

	var p models.Product
	err := r.pool.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.Image,
		product.Price,
	).Scan(
		&p.ProductID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.CreatedAt,
	)
	if err != nil {
		log.Printf("Database error: CreateProduct failed: %v", err)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &p, nil
}

// UpdateProduct updates an existing product
func (r *ProductRepository) UpdateProduct(ctx context.Context, id int, product models.Product) (*models.Product, error) {
	query := `
		UPDATE products
		SET name = $1, description = $2, image = $3, price = $4, "updatedAt" = NOW()
		WHERE "productId" = $5
		RETURNING "productId", name, description, image, price, "createdAt"
	`

	var p models.Product
	err := r.pool.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.Image,
		product.Price,
		id,
	).Scan(
		&p.ProductID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Database error: UpdateProduct(%d) failed: %v", id, err)
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &p, nil
}

// DeleteProduct deletes a product by ID
func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) (*models.Product, error) {
	query := `
		DELETE FROM products 
		WHERE "productId" = $1 
		RETURNING "productId", name, description, image, price, "createdAt"
	`

	var p models.Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ProductID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Database error: DeleteProduct(%d) failed: %v", id, err)
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	return &p, nil
}