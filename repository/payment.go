package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-username/golang-ecommerce-app/models"
)

type Tx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) Create(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	return r.create(ctx, r.pool, payment)
}

func (r *PaymentRepository) CreateWithTx(ctx context.Context, tx Tx, payment models.Payment) (*models.Payment, error) {
	return r.create(ctx, tx, payment)
}

func (r *PaymentRepository) create(ctx context.Context, db Tx, payment models.Payment) (*models.Payment, error) {
	query := `
		INSERT INTO payment ("paymentId", "userId", amount, status, "createdAt")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING "paymentId", "userId", amount, status, "createdAt"
	`

	now := time.Now()
	row := db.QueryRow(ctx, query,
		payment.PaymentID,
		payment.UserId,
		payment.TotalAmount,
		payment.Status,
		now,
	)

	var p models.Payment
	err := row.Scan(
		&p.PaymentID,
		&p.UserId,
		&p.TotalAmount,
		&p.Status,
		&p.CreatedAt,
	)

	if err != nil {
		log.Printf("PaymentRepository.Create failed: %v", err)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &p, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, paymentID string) (*models.Payment, error) {
	query := `
		SELECT "paymentId", "userId", amount, status, "createdAt"
		FROM payment 
		WHERE "paymentId" = $1
	`

	row := r.pool.QueryRow(ctx, query, paymentID)

	var p models.Payment
	err := row.Scan(
		&p.PaymentID,
		&p.UserId,
		&p.TotalAmount,
		&p.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("PaymentRepository.GetByID failed: %v", err)
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &p, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, paymentID, status string) (*models.Payment, error) {
	return r.updateStatus(ctx, r.pool, paymentID, status)
}

func (r *PaymentRepository) UpdateStatusWithTx(ctx context.Context, tx Tx, paymentID, status string) (*models.Payment, error) {
	return r.updateStatus(ctx, tx, paymentID, status)
}

func (r *PaymentRepository) updateStatus(ctx context.Context, db Tx, paymentID, status string) (*models.Payment, error) {
	query := `
		UPDATE payment
		SET status = $1
		WHERE "paymentId" = $2
		RETURNING "paymentId", "userId", amount, status, "createdAt"
	`

	row := db.QueryRow(ctx, query, status, paymentID)

	var p models.Payment
	err := row.Scan(
		&p.PaymentID,
		&p.UserId,
		&p.TotalAmount,
		&p.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Printf("PaymentRepository.UpdateStatus failed: %v", err)
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	return &p, nil
}

func (r *PaymentRepository) ProcessPayment(ctx context.Context, req *models.PaymentRequest, tx Tx) (*models.PaymentResponse, error) {
	// This would typically call the external payment processor
	// For now, we'll simulate a successful payment
	paymentID := fmt.Sprintf("pay_%d", time.Now().UnixNano())

	payment := models.Payment{
		PaymentID:   paymentID,
		UserId:      req.UserID,
		TotalAmount: req.Amount,
		Status:      "success",
	}

	var createdPayment *models.Payment
	var err error

	if tx != nil {
		createdPayment, err = r.CreateWithTx(ctx, tx, payment)
	} else {
		createdPayment, err = r.Create(ctx, payment)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	return &models.PaymentResponse{
		Status:        createdPayment.Status,
		TransactionID: createdPayment.PaymentID,
	}, nil
}
