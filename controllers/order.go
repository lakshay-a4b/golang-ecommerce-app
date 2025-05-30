package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/your-username/golang-ecommerce-app/middlewares"
	"github.com/your-username/golang-ecommerce-app/services"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (oc *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User authentication required")
		return
	}

	createdOrder, err := oc.orderService.CreateOrder(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order created successfully",
		"order":   createdOrder,
	})
}

func (oc *OrderController) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userId, ok := middlewares.GetUserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User authentication required")
		return
	}

	orders, err := oc.orderService.GetUserOrders(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (oc *OrderController) UpdateUserOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["id"]
	if orderId == "" {
		respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}
	status := r.URL.Query().Get("status")
	if status == "" {
		respondWithError(w, http.StatusBadRequest, "Status is required")
		return
	}
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	if orderId == "" {
		respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	_, err := oc.orderService.UpdateUserOrder(r.Context(), userId, orderId, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order updated successfully"})
}
