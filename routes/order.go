package routes

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/controllers"
	"github.com/your-username/golang-ecommerce-app/middlewares"
	"github.com/your-username/golang-ecommerce-app/repository"
	"github.com/your-username/golang-ecommerce-app/services"
)

func RegisterOrderRoutes(r *mux.Router, pool *pgxpool.Pool) {
	orderRepo := repository.NewOrderRepository(pool)
	cartRepo := repository.NewCartRepository(pool)
	paymentRepo := repository.NewPaymentRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	orderService := services.NewOrderService(orderRepo, cartRepo, paymentRepo, productRepo)
	orderController := controllers.NewOrderController(orderService)

	orderRouter := r.PathPrefix("/orders").Subrouter()
	orderRouter.Use(middlewares.AuthenticateToken)

	orderRouter.HandleFunc("/create", orderController.CreateOrder).Methods("POST")
	orderRouter.HandleFunc("/", orderController.GetUserOrders).Methods("GET")
	orderRouter.HandleFunc("/update/{id}", orderController.UpdateUserOrder).Methods("PUT")

}