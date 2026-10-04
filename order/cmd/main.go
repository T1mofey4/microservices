package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	storage := NewOrderStorage()
	inventory := NewInventoryStub()
	payment := &PaymentStub{}
	orderHandler := NewOrderHandler(storage, inventory, payment)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	defer func() {
		if err := server.Close(); err != nil {
			log.Printf("Ошибка закрытия сервера: %v", err)
		}
	}()

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	log.Println("🛑 Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}
	log.Println("✅ Сервер остановлен")
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
		return badRequest("Некорректный UUID заказа"), nil
	}

	order, ok := h.storage.Get(params.OrderUUID)
	if !ok {
		return notFound("Заказ не найден"), nil
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

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req.UserUUID == uuid.Nil {
		return badRequest("Некорректный UUID пользователя"), nil
	}

	if len(req.PartUuids) == 0 {
		return badRequest("Список запрашиваемых запчастей пуст1"), nil
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
