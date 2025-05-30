package routes

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/config"
	"github.com/your-username/golang-ecommerce-app/controllers"
	"github.com/your-username/golang-ecommerce-app/middlewares"
	"github.com/your-username/golang-ecommerce-app/repository"
	"github.com/your-username/golang-ecommerce-app/services"
	"github.com/your-username/golang-ecommerce-app/utils"
)

func RegisterProductRoutes(r *mux.Router, pool *pgxpool.Pool) {
	// Initialize Redis cache
	cache := utils.NewRedisCache(config.RedisClient)
	
	productRepo := repository.NewProductRepository(pool)
	productService := services.NewProductService(productRepo, cache)
	productController := controllers.NewProductController(productService)

	productRouter := r.PathPrefix("/products").Subrouter()

	productRouter.HandleFunc("/", productController.GetAllProducts).Methods("GET")
	productRouter.HandleFunc("/{id}", productController.GetProductById).Methods("GET")

	productAdminRouter := r.PathPrefix("/admin/products").Subrouter()
	productAdminRouter.Use(middlewares.AuthenticateAdminToken) 

	productAdminRouter.HandleFunc("/create", productController.CreateProduct).Methods("POST")
	productAdminRouter.HandleFunc("/update/{id}", productController.UpdateProduct).Methods("PUT")
	productAdminRouter.HandleFunc("/delete/{id}", productController.DeleteProduct).Methods("DELETE")
}