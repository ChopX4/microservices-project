package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderResponse struct {
	OrderUUID  uuid.UUID
	TotalPrice float32
}

type OrderRequest struct {
	UserUUID  uuid.UUID
	PartUUIDs []uuid.UUID
}

type OrderByUUID struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float32
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
	Status          OrderStatus
}

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSPB           PaymentMethod = "SPB"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type PayOrderRequest struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCanceled       OrderStatus = "CANCELED"
)

type Category int32

const (
	CategoryUnknown  Category = 0 // Неизвестная категория
	CategoryEngine   Category = 1 // Двигатель
	CategoryFuel     Category = 2 // Топливо
	CategoryPorthole Category = 3 // Иллюминатор
	CategoryWing     Category = 4 // Крыло
)

type PartsFilter struct {
	UUIDS                   []string
	Names                   []string
	Categories              []Category
	ManunufacturerCountries []string
	Tags                    []string
}

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Manufacturer struct {
	Name    string
	Country string
	WebSite string
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Value interface {
	isKind()
}

type StringValue struct{ V string }
type Int64Value struct{ V int64 }
type Float64Value struct{ V float64 }
type BoolValue struct{ V bool }

func (StringValue) isKind()  {}
func (Int64Value) isKind()   {}
func (Float64Value) isKind() {}
func (BoolValue) isKind()    {}
