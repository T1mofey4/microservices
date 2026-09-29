package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	orderV1 "github.com/T1mofey4/microservices/week_1/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrOrderAlreadyPaid = errors.New("order is already paid")
	ErrOrderCancelled   = errors.New("order cancelled")
	ErrPartNotFound     = errors.New("part not found")
)

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUUIDs       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *orderV1.PaymentMethod
	Status          orderV1.OrderStatus
}

func main() {
	// storage := NewOrderStorage()

	// orderHandler := NewOrderHandler(storage)

	// _, err := orderV1.NewServer(orderHandler)
	// if err != nil {
	// 	log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	// }
}

type PartInfo struct {
	UUID  uuid.UUID
	Price float64
}

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []uuid.UUID) ([]PartInfo, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userUUID, orderUUID uuid.UUID, m orderV1.PaymentMethod) (uuid.UUID, error)
}

type PaymentStub struct{}

func (p *PaymentStub) PayOrder(ctx context.Context, userUUID, orderUUID uuid.UUID, m orderV1.PaymentMethod) (uuid.UUID, error) {
	transactionUUID := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return transactionUUID, nil
}

type InventoryStub struct {
	parts map[uuid.UUID]PartInfo
}

func NewInventoryStub() *InventoryStub {
	id1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	id2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	id3 := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	return &InventoryStub{
		parts: map[uuid.UUID]PartInfo{
			id1: {UUID: id1, Price: 100.0},
			id2: {UUID: id2, Price: 200.0},
			id3: {UUID: id3, Price: 300.0},
		},
	}
}

func (s *InventoryStub) ListParts(ctx context.Context, partUUIDs []uuid.UUID) ([]PartInfo, error) {
	result := make([]PartInfo, 0, len(partUUIDs))

	for _, id := range partUUIDs {
		part, ok := s.parts[id]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrPartNotFound, id)
		}
		result = append(result, part)
	}
	return result, nil
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

	if _, ok := s.orders[order.OrderUUID]; ok {
		return fmt.Errorf("ошибка создания заказа: заказ с UUID %s уже существует", order.OrderUUID)
	}
	snapshot := *order
	s.orders[order.OrderUUID] = &snapshot

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

func (s *OrderStorage) MarkOrderPaid(id, transactionUUID uuid.UUID, m orderV1.PaymentMethod) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return ErrOrderNotFound
	}

	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		switch order.Status {
		case orderV1.OrderStatusPAID:
			return ErrOrderAlreadyPaid
		case orderV1.OrderStatusCANCELLED:
			return ErrOrderCancelled
		default:
			return fmt.Errorf("неизвестный статус: %s", order.Status)
		}
	}

	order.Status = orderV1.OrderStatusPAID
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
	case orderV1.OrderStatusPENDINGPAYMENT:
		order.Status = orderV1.OrderStatusCANCELLED
		return nil
	case orderV1.OrderStatusPAID:
		return ErrOrderAlreadyPaid
	case orderV1.OrderStatusCANCELLED:
		return ErrOrderCancelled
	default:
		return fmt.Errorf("неизвестный статус: %s", order.Status)
	}
}

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

func NewOrderHandler(storage *OrderStorage, inventory InventoryClient, payment PaymentClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventory,
		paymentClient:   payment,
	}
}

func (h *OrderHandler) OrderByUUID(ctx context.Context, params orderV1.OrderByUUIDParams) (orderV1.OrderByUUIDRes, error) {
	if params.OrderUUID == uuid.Nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Некорректный UUID заказа",
		}, nil
	}

	order, ok := h.storage.Get(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Заказ не найден",
		}, nil
	}
	return toOrderDto(order), nil
}

func toOrderDto(order Order) *orderV1.OrderDto {
	dto := &orderV1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptNilUUID(*order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderV1.NewOptNilPaymentMethod(*order.PaymentMethod)
	}

	return dto
}
