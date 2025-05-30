package routes

import (

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/controllers"
	"github.com/your-username/golang-ecommerce-app/repository"
	"github.com/your-username/golang-ecommerce-app/services"
	"github.com/your-username/golang-ecommerce-app/middlewares"
)
func RegisterUserRoutes(r *mux.Router, pool *pgxpool.Pool) {
	userRepo := repository.NewUserRepository(pool)
	userService := services.NewUserService(userRepo)
	controllers := controllers.NewUserController(userService)

	// Public routes
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/signup", controllers.SignupUser).Methods("POST")
	userRouter.HandleFunc("/login", controllers.LoginUser).Methods("POST")

	// Admin subrouter with middleware
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.AuthenticateAdminToken)

	adminRouter.HandleFunc("/user/{userId}", controllers.UpdateUser).Methods("PUT")
	adminRouter.HandleFunc("/user/{userId}", controllers.DeleteUser).Methods("DELETE")

	// SuperAdmin subrouter with middleware
	superAdminRouter := r.PathPrefix("/superadmin").Subrouter()
	superAdminRouter.Use(middlewares.AuthenticateSuperAdminToken)

	superAdminRouter.HandleFunc("/user/role/{userId}", controllers.UpdateUserRole).Methods("PUT")
	superAdminRouter.HandleFunc("/user/{userId}", controllers.DeleteUserAllAccess).Methods("DELETE")
}
