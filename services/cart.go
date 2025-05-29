package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/repository"
)

type CartService struct {
	cartRepo *repository.CartRepository
}

func NewCartService(cartRepo *repository.CartRepository) *CartService {
	return &CartService{cartRepo: cartRepo}
}

func (s *CartService) AddToCartService(ctx context.Context, userID string, newProductInfo json.RawMessage) (*repository.Cart, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if len(newProductInfo) == 0 {
		return nil, errors.New("product info cannot be empty")
	}

	var newProduct models.CartProduct
	if err := json.Unmarshal(newProductInfo, &newProduct); err != nil {
		log.Printf("Invalid product info JSON: %v", err)
		return nil, errors.New("invalid product info")
	}

	if newProduct.ProductID <= 0 || newProduct.Quantity < 1 || newProduct.Price <= 0 {
		return nil, errors.New("invalid product ID, quantity, or price")
	}

	existingCart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error fetching cart for user %s: %v", userID, err)
		return nil, errors.New("failed to get user cart")
	}

	var updatedProducts []models.CartProduct
	if existingCart != nil {
		if err := json.Unmarshal(existingCart.ProductInfo, &updatedProducts); err != nil {
			log.Printf("Failed to parse existing product info: %v", err)
			return nil, errors.New("failed to process existing cart")
		}

		found := false
		for i, p := range updatedProducts {
			if p.ProductID == newProduct.ProductID {
				updatedProducts[i].Quantity += newProduct.Quantity
				updatedProducts[i].Price = newProduct.Price
				found = true
				break
			}
		}
		if !found {
			updatedProducts = append(updatedProducts, newProduct)
		}
	} else {
		updatedProducts = []models.CartProduct{newProduct}
	}

	updatedJson, err := json.Marshal(updatedProducts)
	if err != nil {
		log.Printf("Failed to marshal updated cart: %v", err)
		return nil, errors.New("failed to update cart")
	}

	return s.cartRepo.CreateOrUpdateCart(ctx, userID, updatedJson)
}

func (s *CartService) UpdateCartService(ctx context.Context, userID string, productInfo json.RawMessage) (*repository.Cart, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if len(productInfo) == 0 {
		return nil, errors.New("product info cannot be empty")
	}

	// Validate product info structure
	var products []models.CartProduct
	if err := json.Unmarshal(productInfo, &products); err != nil {
		return nil, errors.New("invalid product info format")
	}

	for _, p := range products {
		if p.ProductID <= 0 || p.Quantity < 1 || p.Price <= 0 {
			return nil, errors.New("invalid product data in cart")
		}
	}

	return s.cartRepo.UpdateCartProducts(ctx, userID, productInfo)
}

func (s *CartService) RemoveFromCartService(ctx context.Context, userID string, productID int, quantityToRemove int) (*repository.Cart, error) {
    if userID == "" {
        return nil, errors.New("user ID cannot be empty")
    }
    if productID <= 0 {
        return nil, errors.New("invalid product ID")
    }
    if quantityToRemove <= 0 {
        return nil, errors.New("quantity to remove must be positive")
    }

    existingCart, err := s.cartRepo.GetCartByUserID(ctx, userID)
    if err != nil {
        log.Printf("Error fetching cart for user %s: %v", userID, err)
        return nil, errors.New("failed to get user cart")
    }
    if existingCart == nil {
        return nil, errors.New("cart not found")
    }

    var products []models.CartProduct
    if err := json.Unmarshal(existingCart.ProductInfo, &products); err != nil {
        log.Printf("Failed to parse existing product info: %v", err)
        return nil, errors.New("failed to process existing cart")
    }

    updatedProducts := make([]models.CartProduct, 0, len(products))
    productFound := false

    for _, p := range products {
        if p.ProductID == productID {
            productFound = true
            newQuantity := p.Quantity - quantityToRemove
            if newQuantity > 0 {
                // Keep the product with reduced quantity
                updatedProducts = append(updatedProducts, models.CartProduct{
                    ProductID: p.ProductID,
                    Quantity:  newQuantity,
                    Price:     p.Price,
                })
            }
            // If newQuantity <= 0, we don't add it back (effectively removing it)
        } else {
            updatedProducts = append(updatedProducts, p)
        }
    }

    if !productFound {
        return nil, errors.New("product not found in cart")
    }

    updatedJson, err := json.Marshal(updatedProducts)
    if err != nil {
        log.Printf("Failed to marshal updated cart: %v", err)
        return nil, errors.New("failed to update cart")
    }

    return s.cartRepo.UpdateCartProducts(ctx, userID, updatedJson)
}

func (s *CartService) GetCartService(ctx context.Context, userID string) (*repository.Cart, error) {
	return s.cartRepo.GetCartByUserID(ctx, userID)
}

func (s *CartService) ClearUserCart(ctx context.Context, userID string) error {
	return s.cartRepo.DeleteCart(ctx, userID, nil)
}