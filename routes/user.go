package routes

import (

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/controllers"
	"github.com/your-username/golang-ecommerce-app/repository"
	"github.com/your-username/golang-ecommerce-app/services"
)

func RegisterUserRoutes(r *mux.Router , pool *pgxpool.Pool) {
	
	userRepo := repository.NewUserRepository(pool)
	userService := services.NewUserService(userRepo)
	controllers := controllers.NewUserController(userService)

	userRouter := r.PathPrefix("/users").Subrouter()
	
	userRouter.HandleFunc("/signup", controllers.SignupUser).Methods("POST")
	userRouter.HandleFunc("/login", controllers.LoginUser).Methods("POST")

}