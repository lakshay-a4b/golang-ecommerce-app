package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/repository"
	"github.com/your-username/golang-ecommerce-app/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type ServiceError struct {
	Status  int
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (u *UserService) Signup(ctx context.Context, userData models.User) (*models.User, error) {
	existingUser, err := u.userRepo.FindByuserId(ctx, userData.UserId)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}

	if existingUser != nil {
		return nil, &ServiceError{
			Status:  400,
			Message: "User already exists",
		}
	}

	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		log.Println("Hashing error:", err)
		return nil, err
	}
	userData.Password = hashedPassword
	result, err := u.userRepo.CreateUser(ctx, userData)

	if err != nil {
		return nil, err
	}

	event := models.UserEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Name:      userData.UserId,
		Email:     userData.Email,
		Role:      userData.Role,
		Action:    "signup",
		UserId:    userData.UserId,
	}

	eventMap := map[string]interface{}{
		"Timestamp": event.Timestamp,
		"Name":      event.Name,
		"Email":     event.Email,
		"Role":      event.Role,
		"Action":    event.Action,
		"userId":    event.UserId,
	}

	go utils.LogEventToProducer("User Signup", userData.UserId, eventMap)

	return result, nil
}

func (u *UserService) Login(ctx context.Context, userId, password string) (string, error) {
	user, err := u.userRepo.FindByuserId(ctx, userId)
	if err != nil {
		return "", &ServiceError{
			Status:  401,
			Message: "Invalid credentials",
		}
	}
	if user == nil {
		return "", &ServiceError{
			Status:  401,
			Message: "Invalid credentials",
		}
	}

	if !utils.ComparePasswords(password, user.Password) {
		return "", &ServiceError{
			Status:  401,
			Message: "Invalid credentials",
		}
	}

	token, err := utils.GenerateToken(user.UserId, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) UpdateUserService(ctx context.Context, userId string, updates map[string]interface{}) (*models.User, error) {
	user, err := u.userRepo.FindByuserId(ctx, userId)
	if err != nil || user == nil {
		return nil, &ServiceError{
			Status:  404,
			Message: "User not found",
		}
	}

	updatedUser, err := u.userRepo.UpdateUser(ctx, userId, updates)
	if err != nil {
		log.Println("Error updating user:", err)
		return nil, err
	}

	return updatedUser, nil
}

func (u *UserService) DeleteUserService(ctx context.Context, userId string) error {
	user, err := u.userRepo.FindByuserId(ctx, userId)
	if err != nil || user == nil {
		return &ServiceError{
			Status:  404,
			Message: "User not found",
		}
	}

	if _, err := u.userRepo.DeleteUser(ctx, userId); err != nil {
		log.Println("Error deleting user:", err)
		return err
	}

	return nil
}
