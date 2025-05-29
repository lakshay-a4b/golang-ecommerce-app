package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/your-username/golang-ecommerce-app/middlewares"
	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/services"
)

type CartController struct {
	cartService *services.CartService
}

func NewCartController(cartService *services.CartService) *CartController {
	return &CartController{cartService: cartService}
}

func (cc *CartController) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User authentication required")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	var product models.CartProduct
	if err := json.Unmarshal(bodyBytes, &product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON structure")
		return
	}

	if product.ProductID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}
	if product.Quantity <= 0 {
		respondWithError(w, http.StatusBadRequest, "Quantity must be positive")
		return
	}
	if product.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "Price must be positive")
		return
	}

	cart, err := cc.cartService.AddToCartService(r.Context(), userID, bodyBytes)
	if err != nil {
		log.Printf("Error adding to cart: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to add to cart")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"cart":    cart,
	})
}

func (cc *CartController) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User authentication required")
		return
	}

	cart, err := cc.cartService.GetCartService(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting cart: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch cart")
		return
	}

	if cart == nil {
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"cart":    nil,
			"message": "Cart is empty",
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"cart":    cart,
	})
}

func (cc *CartController) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
    userID, ok := middlewares.GetUserFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "User authentication required")
        return
    }

    vars := mux.Vars(r)
    productIDStr := vars["productId"]
    if productIDStr == "" {
        respondWithError(w, http.StatusBadRequest, "Product ID is required")
        return
    }

    productID, err := strconv.Atoi(productIDStr)
    if err != nil || productID <= 0 {
        respondWithError(w, http.StatusBadRequest, "Invalid product ID")
        return
    }

    // Parse request body for quantity to remove
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    defer r.Body.Close()

    var requestBody struct {
        Quantity int `json:"quantity"`
    }
    if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid JSON structure")
        return
    }

    if requestBody.Quantity <= 0 {
        respondWithError(w, http.StatusBadRequest, "Quantity to remove must be positive")
        return
    }

    updatedCart, err := cc.cartService.RemoveFromCartService(r.Context(), userID, productID, requestBody.Quantity)
    if err != nil {
        log.Printf("Error removing from cart: %v", err)
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]interface{}{
        "success": true,
        "message": "Product quantity reduced in cart",
        "cart":    updatedCart,
    })
}

func (cc *CartController) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User authentication required")
		return
	}

	err := cc.cartService.ClearUserCart(r.Context(), userID)
	if err != nil {
		log.Printf("Error clearing cart: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to clear cart")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Cart cleared successfully",
	})
}

func (cc *CartController) UpdateCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User authentication required")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	cart, err := cc.cartService.UpdateCartService(r.Context(), userID, bodyBytes)
	if err != nil {
		log.Printf("Error updating cart: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update cart")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"cart":    cart,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}