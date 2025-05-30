package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/your-username/golang-ecommerce-app/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header format must be 'Bearer {token}'")
			return
		}

		tokenString := parts[1]
		user, err := utils.VerifyToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateToken for admin role
func AuthenticateAdminToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header format must be 'Bearer {token}'")
			return
		}

		tokenString := parts[1]

		// Verify token with admin or superadmin role
		user, err := utils.VerifyTokenWithRoles(tokenString, []string{"admin", "superadmin"})

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Not Authorized: "+err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateToken for super admin role
func AuthenticateSuperAdminToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header format must be 'Bearer {token}'")
			return
		}

		tokenString := parts[1]
		user, err := utils.VerifyTokenWithRoles(tokenString, []string{"superadmin"})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Not Authorized: "+err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(UserContextKey).(string)
	return user, ok
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   true,
		"message": message,
	})
}