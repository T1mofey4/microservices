package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"uuid"

	orderV1 "github.com/T1mofey4/microservices/week_1/shared/pkg/openapi/order/v1"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrOrderAlreadyPaid = errors.New("order is already paid")
	ErrOrderCancelled   = errors.New("order cancelled")
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
}

func main() {
	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	_, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*Order),
	}
}

func (s *OrderStorage) Create(order *Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	copyOrder := *order

	if _, ok := s.orders[order.OrderUUID]; ok {
		return fmt.Errorf("ошибка создания заказа: заказ с UUID %s уже существует", order.OrderUUID)
	}
	s.orders[order.OrderUUID] = &copyOrder

	return nil
}

func (s *OrderStorage) Get(orderUUID uuid.UUID) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return Order{}, false
	}

	return *order, true
}

func (s *OrderStorage) MarkOrderPaid(id, transactionUUID uuid.UUID, m PaymentMethod) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return ErrOrderNotFound
	}

	if order.Status != OrderStatusPendingPayment {
		switch order.Status {
		case OrderStatusPaid:
			return ErrOrderAlreadyPaid
		case OrderStatusCancelled:
			return ErrOrderCancelled
		default:
			return fmt.Errorf("неизвестный статус: %s", order.Status)
		}
	}

	order.Status = OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &m

	return nil
}

func (s *OrderStorage) MarkOrderCancelled(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return ErrOrderNotFound
	}

	switch order.Status {
	case OrderStatusPendingPayment:
		order.Status = OrderStatusCancelled
		return nil
	case OrderStatusPaid:
		return ErrOrderAlreadyPaid
	case OrderStatusCancelled:
		return ErrOrderCancelled
	default:
		return fmt.Errorf("неизвестный статус: %s", order.Status)
	}
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}
