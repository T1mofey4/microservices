package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
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

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*Order),
	}
}

// Записать информацию о заказе в storage
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

// Получить информацию о заказе из storage
func (s *OrderStorage) Get(orderUUID uuid.UUID) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return Order{}, false
	}

	return *order, true
}

// Проверяет и меняет статус заказа
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

// Меняет статус заказа на Cancelled
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

// Получить заказ по uuid
func (h *OrderHandler) OrderByUUID(ctx context.Context, params orderV1.OrderByUUIDParams) (orderV1.OrderByUUIDRes, error) {
	if params.OrderUUID == uuid.Nil {
		return badRequest("Некорректный UUID заказа"), nil
	}

	order, ok := h.storage.Get(params.OrderUUID)
	if !ok {
		return notFound("Заказ не найден"), nil
	}
	return toOrderDto(order), nil
}

// Преобразование Order в OrderDto для ogen
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

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req.UserUUID == uuid.Nil {
		return badRequest("Некорректный UUID пользователя"), nil
	}

	if len(req.PartUuids) == 0 {
		return badRequest("Список запрашиваемых запчастей пуст"), nil
	}

	parts, err := h.inventoryClient.ListParts(ctx, req.PartUuids)
	if err != nil {
		if errors.Is(err, ErrPartNotFound) {
			return &orderV1.BadRequestError{Code: 400, Message: "одна или несколько деталей не найдены"}, nil
		}
		log.Printf("ListParts failed: %v", err)
		return internalError("внутренняя ошибка"), nil
	}

	var total float64
	for _, part := range parts {
		total += part.Price
	}

	order := &Order{
		OrderUUID:  uuid.New(),
		UserUUID:   req.UserUUID,
		PartUUIDs:  req.PartUuids,
		TotalPrice: total,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	if err := h.storage.Create(order); err != nil {
		log.Printf("Create order failed: %v", err)
		return internalError("Ошибка сохранения заказа"), nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (h *OrderHandler) OrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPayParams) (orderV1.OrderPayRes, error) {
	if params.OrderUUID == uuid.Nil {
		return badRequest("Некорректный UUID заказа"), nil
	}

	order, ok := h.storage.Get(params.OrderUUID)
	if !ok {
		return notFound("Заказ не найден"), nil
	}

	switch order.Status {
	case orderV1.OrderStatusPENDINGPAYMENT:
	case orderV1.OrderStatusPAID:
		return badRequest("Заказ уже оплачен"), nil
	case orderV1.OrderStatusCANCELLED:
		return badRequest("Нельзя оплатить отмененный заказ"), nil
	default:
		return internalError("Неизвестный статус заказа"), nil
	}

	tUUID, err := h.paymentClient.PayOrder(ctx, order.UserUUID, order.OrderUUID, req.PaymentMethod)
	if err != nil {
		log.Printf("PayOrder failed: %v", err)
		return internalError("Не удалось оплатить заказ"), nil
	}

	if err := h.storage.MarkOrderPaid(order.OrderUUID, tUUID, req.PaymentMethod); err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound):
			return notFound("Заказ не найден"), nil
		case errors.Is(err, ErrOrderAlreadyPaid):
			return badRequest("Заказ уже оплачен"), nil
		case errors.Is(err, ErrOrderCancelled):
			return badRequest("Нельзя оплатить отмененный заказ"), nil
		default:
			log.Printf("MarkOrderPaid failed: %v", err)
			return internalError("Внутренняя ошибка"), nil
		}
	}

	return &orderV1.PayOrderResponse{TransactionUUID: tUUID}, nil
}

func (h *OrderHandler) OrderCancel(ctx context.Context, params orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	err := h.storage.MarkOrderCancelled(params.OrderUUID)
	switch {
	case err == nil:
		return &orderV1.OrderCancelNoContent{}, nil
	case errors.Is(err, ErrOrderNotFound):
		return notFound("Заказ не найден"), nil
	case errors.Is(err, ErrOrderAlreadyPaid):
		return conflictError("Заказ уже оплачен и не может быть отменён"), nil
	case errors.Is(err, ErrOrderCancelled):
		return conflictError("Заказ уже отменен"), nil
	default:
		log.Printf("OrderCancel failed: %v", err)
		return internalError("Внутренняя ошибка"), nil
	}
}

func badRequest(msg string) *orderV1.BadRequestError {
	return &orderV1.BadRequestError{Code: 400, Message: msg}
}

func notFound(msg string) *orderV1.NotFoundError {
	return &orderV1.NotFoundError{Code: 404, Message: msg}
}

func internalError(msg string) *orderV1.InternalServerError {
	return &orderV1.InternalServerError{Code: 500, Message: msg}
}

func conflictError(msg string) *orderV1.ConflictError {
	return &orderV1.ConflictError{Code: 409, Message: msg}
}
