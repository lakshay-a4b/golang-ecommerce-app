package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/models"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// FindByuserId fetches user by ID
func (r *UserRepository) FindByuserId(ctx context.Context, userId string) (*models.User, error) {
	query := `SELECT "userId", email, password, role, "createdAt" FROM users WHERE "userId" = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, userId).
		Scan(&user.UserId, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error finding user by ID: %v", err)
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &user, nil
}

// CreateUser inserts a new user
func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	query := `
		INSERT INTO users ("userId", email, password, role, "createdAt")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING "userId", email, password, role, "createdAt"
	`

	var newUser models.User
	err := r.pool.QueryRow(ctx, query,
		user.UserId,
		user.Email,
		user.Password,
		user.Role,
		user.CreatedAt,
	).Scan(
		&newUser.UserId,
		&newUser.Email,
		&newUser.Password,
		&newUser.Role,
		&newUser.CreatedAt,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

// FindByEmail fetches user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT "userId", email, password, role, "createdAt" FROM users WHERE email = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).
		Scan(&user.UserId, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error finding user by email: %v", err)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

// UpdateUser updates user details
func (r *UserRepository) UpdateUser(ctx context.Context, userId string, updates map[string]interface{}) (*models.User, error) {
	query := `
		UPDATE users
		SET email = COALESCE($1, email),
			password = COALESCE($2, password),
			"updatedAt" = NOW()
		WHERE "userId" = $3
		RETURNING "userId", email, password, role, "createdAt"
	`

	var email *string
	var password *string

	if val, ok := updates["email"]; ok {
		if s, ok := val.(string); ok {
			email = &s
		}
	}

	if val, ok := updates["password"]; ok {
		if s, ok := val.(string); ok {
			password = &s
		}
	}

	var updatedUser models.User
	err := r.pool.QueryRow(ctx, query, email, password, userId).
		Scan(
			&updatedUser.UserId,
			&updatedUser.Email,
			&updatedUser.Password,
			&updatedUser.Role,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error updating user: %v", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedUser, nil
}

// DeleteUser deletes user by ID
func (r *UserRepository) DeleteUser(ctx context.Context, userId string) (*models.User, error) {
	query := `
		DELETE FROM users 
		WHERE "userId" = $1 
		RETURNING "userId", email, password, role, "createdAt"
	`

	var deletedUser models.User
	err := r.pool.QueryRow(ctx, query, userId).
		Scan(
			&deletedUser.UserId,
			&deletedUser.Email,
			&deletedUser.Password,
			&deletedUser.Role,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Error deleting user: %v", err)
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &deletedUser, nil
}
