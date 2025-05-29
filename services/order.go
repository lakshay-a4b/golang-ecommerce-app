package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	paymentRepo *repository.PaymentRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	cartRepo *repository.CartRepository,
	paymentRepo *repository.PaymentRepository,
	productRepo *repository.ProductRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		paymentRepo: paymentRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userId string) (*models.Order, error) {
	if userId == "" {
		return nil, fmt.Errorf("invalid user ID")
	}

	cart, err := s.cartRepo.GetCartByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	if cart == nil || len(cart.ProductInfo) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	var cartProducts []models.CartProduct
	if err := json.Unmarshal(cart.ProductInfo, &cartProducts); err != nil {
		return nil, fmt.Errorf("failed to parse cart products: %w", err)
	}

	var items []models.OrderItem
	var totalAmount float64

	for _, item := range cartProducts {
		product, err := s.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}
		if product == nil {
			return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}

		items = append(items, models.OrderItem{
			ProductID: product.ProductID,
			Name:      product.Name,
			Image:     product.Image,
			Price:     product.Price,
			Quantity:  item.Quantity,
		})

		totalAmount += product.Price * float64(item.Quantity)
	}

	// tx, err := s.orderRepo.BeginTx(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to begin transaction: %w", err)
	// }
	// defer func() {
	// 	if err != nil {
	// 		if rbErr := tx.Rollback(ctx); rbErr != nil {
	// 			log.Printf("Failed to rollback transaction: %v", rbErr)
	// 		}
	// 	}
	// }()

	paymentResult, err := s.paymentRepo.ProcessPayment(ctx, &models.PaymentRequest{
		UserID: userId,
		Amount: totalAmount,
	},nil)
	if err != nil {
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}
	if paymentResult == nil || paymentResult.TransactionID == "" {
		return nil, fmt.Errorf("payment processing failed")
	}

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order items: %w", err)
	}

	order := &models.Order{
		PaymentID:   paymentResult.TransactionID,
		UserId:      userId,
		ProductInfo: itemsJSON,
		Status:      "Order-Accepted",
	}

	createdOrder, err := s.orderRepo.AddOrder(ctx, *order, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	if err := s.cartRepo.DeleteCart(ctx, userId, nil); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// if err := tx.Commit(ctx); err != nil {
	// 	return nil, fmt.Errorf("failed to commit transaction: %w", err)
	// }

	createdOrder.ProductInfo, err = json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order items for response: %w", err)
	}
	return createdOrder, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userId string) ([]models.Order, error) {
	orders, err := s.orderRepo.GetOrdersByUser(ctx, userId)
	if err != nil {
		log.Printf("Error fetching orders for user %s: %v", userId, err)
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (s *OrderService) UpdateUserOrder(ctx context.Context, userId string, orderId string, status string) (*models.Order, error) {
	if userId == "" || orderId == "" {
		return nil, fmt.Errorf("invalid user ID or order ID")
	}

	order, err := s.orderRepo.GetOrderByID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil || order.UserId != userId {
		return nil, fmt.Errorf("order not found or does not belong to user")
	}

	order.Status = status
	updatedOrder, err := s.orderRepo.UpdateOrder(ctx, *order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return updatedOrder, nil
}
func (s *OrderService) GetOrderByID(ctx context.Context, orderId string) (*models.Order, error) {
	if orderId == "" {
		return nil, fmt.Errorf("invalid order ID")
	}

	order, err := s.orderRepo.GetOrderByID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}