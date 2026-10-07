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
	"syscall"
	"time"

	orderV1 "github.com/T1mofey4/microservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/T1mofey4/microservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/T1mofey4/microservices/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort          = "8080"
	inventoryPort     = "50051"
	paymentPort       = "50052"
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
	inventoryAddr := net.JoinHostPort("localhost", inventoryPort)
	inventoryConn, err := grpc.NewClient(
		inventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("ошибка создания клиента Inventory: %v", err)
	}
	inventoryClient := NewInventoryGRPCClient(inventoryConn)
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения Inventory: %v", err)
		}
	}()

	paymentAddr := net.JoinHostPort("localhost", paymentPort)
	paymentConn, err := grpc.NewClient(
		paymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("ошибка создания клиента Payment: %v", err)
	}
	paymentClient := NewPaymentGRPCClient(paymentConn)
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("Ошибка закрытия соединения Payment: %v", err)
		}
	}()

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)
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
		err := server.ListenAndServe()
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

type InventoryGRPCClient struct {
	client inventoryV1.InventoryServiceClient
}

func NewInventoryGRPCClient(conn *grpc.ClientConn) *InventoryGRPCClient {
	return &InventoryGRPCClient{client: inventoryV1.NewInventoryServiceClient(conn)}
}

func (c *InventoryGRPCClient) ListParts(ctx context.Context, partUUIDs []uuid.UUID) ([]PartInfo, error) {
	partUUID := make([]string, 0, len(partUUIDs))
	for _, u := range partUUIDs {
		partUUID = append(partUUID, u.String())
	}

	filter := &inventoryV1.PartsFilter{Uuids: partUUID}
	req := &inventoryV1.ListPartsRequest{Filter: filter}
	res, err := c.client.ListParts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("inventory ListParts: %w", err)
	}

	if len(res.Parts) != len(partUUID) {
		return nil, ErrPartNotFound
	}

	result := make([]PartInfo, 0, len(res.Parts))
	for _, p := range res.Parts {
		uuidString, err := uuid.Parse(p.Uuid)
		if err != nil {
			return nil, err
		}

		toPartInfo := PartInfo{
			UUID:  uuidString,
			Price: p.Price,
		}
		result = append(result, toPartInfo)
	}

	return result, nil
}

type PaymentGRPCClient struct {
	client paymentV1.PaymentServiceClient
}

func NewPaymentGRPCClient(conn *grpc.ClientConn) *PaymentGRPCClient {
	return &PaymentGRPCClient{client: paymentV1.NewPaymentServiceClient(conn)}
}

func (c *PaymentGRPCClient) PayOrder(ctx context.Context, userUUID, orderUUID uuid.UUID, m orderV1.PaymentMethod) (uuid.UUID, error) {
	mapped, err := mapPaymentMethod(m)
	if err != nil {
		return uuid.Nil, err
	}

	req := &paymentV1.PayOrderRequest{
		UserUuid:      userUUID.String(),
		OrderUuid:     orderUUID.String(),
		PaymentMethod: mapped,
	}

	res, err := c.client.PayOrder(ctx, req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка запроса на оплату: %w", err)
	}

	transactionUUID, err := uuid.Parse(res.TransactionUuid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка преобразования uuid: %w", err)
	}

	return transactionUUID, nil
}

func mapPaymentMethod(m orderV1.PaymentMethod) (paymentV1.PaymentMethod, error) {
	switch m {
	case orderV1.PaymentMethodUNKNOWN:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, errors.New("способ оплаты не выбран")
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, fmt.Errorf("неизвестный метод оплаты: %s", m)
	}
}
