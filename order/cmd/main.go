package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	order_v1 "github.com/ChopX4/raketka/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort         = "8080"
	paymentAddress   = "localhost:50052"
	inventoryAddress = "localhost:50051"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*order_v1.OrderByUUID
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*order_v1.OrderByUUID),
	}
}

func (s *OrderStorage) CreateOrder(order *order_v1.OrderByUUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order.Status = order_v1.OrderStatusPENDINGPAYMENT
	s.orders[order.OrderUUID.String()] = order

	return nil
}

func (s *OrderStorage) PayOrder(orderUUID, transactionUUID string, paymentMethod order_v1.PaymentMethod) (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return uuid.Nil, fmt.Errorf("order %s not found", orderUUID)
	}

	UUIDtransactionUUID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse payment transaction id: %w", err)
	}

	order.Status = order_v1.OrderStatusPAID
	order.TransactionUUID = UUIDtransactionUUID
	order.PaymentMethod = paymentMethod
	return UUIDtransactionUUID, nil
}

func (s *OrderStorage) GetOrder(orderUUID string) (*order_v1.OrderByUUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return nil, fmt.Errorf("order %s not found", orderUUID)
	}

	return order, nil
}

func (s *OrderStorage) CancelOrder(orderUUID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderUUID]
	if !ok {
		return fmt.Errorf("not found")
	}

	if order.Status == order_v1.OrderStatusCANCELED {
		return fmt.Errorf("conflict")
	}

	order.Status = order_v1.OrderStatusCANCELED
	return nil
}

type OrderHandler struct {
	storage         *OrderStorage
	paymentClient   payment_v1.PaymentServiceClient
	inventoryClient inventory_v1.InventoryServiceClient
}

func NewOrderHandler(storage *OrderStorage, paymentClient payment_v1.PaymentServiceClient, inventoryClient inventory_v1.InventoryServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}

func (h *OrderHandler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

func (h *OrderHandler) CancelOrder(_ context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	err := h.storage.CancelOrder(params.OrderUUID.String())
	if err != nil {
		switch err.Error() {
		case "not_found":
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		case "conflict":
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Conflict",
			}, nil
		default:
			return nil, err
		}
	}

	return &order_v1.NoContentError{
		Code:    204,
		Message: "Order successfully canceled",
	}, nil
}

func (h *OrderHandler) CreateOrder(_ context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	uuids := make([]string, 0, len(req.PartUuids))
	for _, v := range req.PartUuids {
		uuids = append(uuids, v.String())
	}

	grpcReq := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: uuids,
		},
	}

	responseParts, err := h.inventoryClient.ListParts(context.Background(), grpcReq)
	if err != nil {
		return nil, err
	}

	if len(responseParts.GetParts()) != len(req.PartUuids) {
		return &order_v1.BadRequestError{
			Code:    400,
			Message: "Bad request - validation error",
		}, nil
	}

	parts := responseParts.GetParts()

	var totalPrice float64
	partsUUIDS := make([]uuid.UUID, 0, len(parts))
	orderUUID := uuid.New()

	for _, v := range parts {
		totalPrice += v.GetPrice()
		uuidPart, err := uuid.Parse(v.GetUuid())
		if err != nil {
			return &order_v1.BadRequestError{
				Code:    400,
				Message: "Bad request - validation error",
			}, nil
		}
		partsUUIDS = append(partsUUIDS, uuidPart)
	}

	if err := h.storage.CreateOrder(&order_v1.OrderByUUID{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  partsUUIDS,
		TotalPrice: float32(totalPrice),
	}); err != nil {
		return nil, err
	}

	return &order_v1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: float32(totalPrice),
	}, nil
}

func (h *OrderHandler) GetOrderById(_ context.Context, params order_v1.GetOrderByIdParams) (order_v1.GetOrderByIdRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) OrderPay(_ context.Context, req *order_v1.OrderPayReq, params order_v1.OrderPayParams) (order_v1.OrderPayRes, error) {
	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	method := req.GetPaymentMethod()
	var grpcMethod payment_v1.PaymentMethod

	switch method.Value {
	case order_v1.PaymentMethodCARD:
		grpcMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case order_v1.PaymentMethodSPB:
		grpcMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_SPB
	case order_v1.PaymentMethodCREDITCARD:
		grpcMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case order_v1.PaymentMethodINVESTORMONEY:
		grpcMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	}

	grpcRequest := &payment_v1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: grpcMethod,
	}

	paymentResponse, err := h.paymentClient.PayOrder(context.Background(), grpcRequest)
	if err != nil {
		return nil, err
	}

	transactionUUID, err := h.storage.PayOrder(params.OrderUUID.String(), paymentResponse.GetTransactionUuid(), method.Value)
	if err != nil {
		return &order_v1.BadRequestError{
			Code:    400,
			Message: "Bad request - validation error",
		}, nil
	}

	return &order_v1.TransactionUUID{
		TransactionUUID: transactionUUID,
	}, nil
}

func main() {
	storage := NewOrderStorage()

	paymentCon, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentCon.Close(); cerr != nil {
			log.Fatalf("failed to close connect: %v", cerr)
		}
	}()

	inventoryCon, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryCon.Close(); cerr != nil {
			log.Fatalf("failed to close connect: %v", cerr)
		}
	}()

	paymentClient := payment_v1.NewPaymentServiceClient(paymentCon)
	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryCon)

	orderHandler := NewOrderHandler(storage, paymentClient, inventoryClient)

	orderServer, err := order_v1.NewServer(orderHandler)
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

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
