package routes

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/controllers"
	"github.com/your-username/golang-ecommerce-app/middlewares"
	"github.com/your-username/golang-ecommerce-app/repository"
	"github.com/your-username/golang-ecommerce-app/services"
)

func RegisterCartRoutes(r *mux.Router, pool *pgxpool.Pool) {
	// Initialize dependencies
	cartRepo := repository.NewCartRepository(pool)
	cartService := services.NewCartService(cartRepo)
	cartController := controllers.NewCartController(cartService)

	cartRouter := r.PathPrefix("/cart").Subrouter()
	cartRouter.Use(middlewares.AuthenticateToken)

	// RESTful routes
	cartRouter.HandleFunc("/", cartController.GetCart).Methods("GET")
	cartRouter.HandleFunc("/add", cartController.AddToCart).Methods("POST")
	cartRouter.HandleFunc("/{productId}", cartController.RemoveFromCart).Methods("DELETE")
}