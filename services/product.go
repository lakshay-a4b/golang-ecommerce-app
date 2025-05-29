package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/repository"
	"github.com/your-username/golang-ecommerce-app/utils"
)

const (
	productCacheTTL       = time.Hour
	defaultPage           = 1
	defaultLimit          = 5
	firstPageCachePattern = "product:1:*" // Pattern to match first page cache keys
)

type ProductService struct {
	productRepo *repository.ProductRepository
	cache       utils.CacheProvider
}

func NewProductService(productRepo *repository.ProductRepository, cache utils.CacheProvider) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		cache:       cache,
	}
}

func (s *ProductService) GetPaginatedProducts(ctx context.Context, page, limit int) (*models.PaginatedProductResponse, error) {
	// Validate and set defaults
	if page <= 0 {
		page = defaultPage
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	offset := (page - 1) * limit

	cacheKey := fmt.Sprintf("products:%d:%d", page, limit)

	var cachedResponse models.PaginatedProductResponse
	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		log.Printf("Cache hit for key: %s", cacheKey)
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			return &cachedResponse, nil
		}
	}

	products, err := s.productRepo.GetPaginatedProducts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get paginated products: %w", err)
	}

	total, err := s.productRepo.GetTotalProductCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get product count: %w", err)
	}

	response := &models.PaginatedProductResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal products for caching: %v", err)
	} else {
		if err := s.cache.Set(ctx, cacheKey, string(jsonData), productCacheTTL); err != nil {
			log.Printf("Failed to cache products: %v", err)
		}
	}

	return response, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}

	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}
	if product == nil {
		return nil, nil
	}

	return product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	if product.Name == "" {
		return nil, errors.New("product name is required")
	}
	if product.Price <= 0 {
		return nil, errors.New("product price must be positive")
	}

	createdProduct, err := s.productRepo.CreateProduct(ctx, *product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	if err := s.cache.DeletePattern(ctx, firstPageCachePattern); err != nil {
		log.Printf("Failed to invalidate product cache: %v", err)
	}

	return createdProduct, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, updates models.Product) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}

	if updates.Name == "" && updates.Description == "" && updates.Image == "" && updates.Price == 0 {
		return nil, errors.New("no valid fields provided for update")
	}

	updatedProduct, err := s.productRepo.UpdateProduct(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	if updatedProduct == nil {
		return nil, nil
	}

	if err := s.cache.DeletePattern(ctx, "products:*"); err != nil {
		log.Printf("Failed to invalidate product cache: %v", err)
	}

	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) (*models.Product, error) {
	if id <= 0 {
		return nil ,errors.New("invalid product ID")
	}

	deletedProduct, err := s.productRepo.DeleteProduct(ctx, id)
	if err != nil {
		return deletedProduct, fmt.Errorf("failed to delete product: %w", err)
	}
	if deletedProduct == nil {
		return nil, nil
	}

	if err := s.cache.DeletePattern(ctx, "products:*"); err != nil {
		log.Printf("Failed to invalidate product cache: %v", err)
	}

	return deletedProduct, nil
}