package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/repository"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
	}
}

type PaymentServiceInterface interface {
	ProcessPayment(ctx context.Context, orderDetails *models.PaymentRequest, tx repository.Tx) (*models.PaymentResponse, error)
}

func (s *PaymentService) ProcessPayment(ctx context.Context, paymentRequest *models.PaymentRequest, tx pgx.Tx) (*models.PaymentResponse, error) {
	if paymentRequest.UserID == "" {
		return nil, errors.New("user ID is required")
	}
	if paymentRequest.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	paymentID := fmt.Sprintf("pay_%d", time.Now().UnixMilli())

	payment := models.Payment{
		PaymentID:   paymentID,
		UserId:      paymentRequest.UserID,
		TotalAmount: paymentRequest.Amount,
		Status:      "pending",
	}

	var createdPayment *models.Payment
	var err error

	if tx != nil {
		createdPayment, err = s.paymentRepo.CreateWithTx(ctx, tx, payment)
	} else {
		createdPayment, err = s.paymentRepo.Create(ctx, payment)
	}

	if err != nil {
		log.Printf("PaymentService.ProcessPayment failed: %v", err)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Simulate payment processing
	createdPayment.Status = "success"
	if tx != nil {
		_, err = s.paymentRepo.UpdateStatusWithTx(ctx, tx, paymentID, "success")
	} else {
		_, err = s.paymentRepo.UpdateStatus(ctx, paymentID, "success")
	}

	if err != nil {
		log.Printf("PaymentService.UpdateStatus failed: %v", err)
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	return &models.PaymentResponse{
		Status:        createdPayment.Status,
		TransactionID: createdPayment.PaymentID,
	}, nil
}