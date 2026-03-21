package model

import "github.com/google/uuid"

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

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCanceled       OrderStatus = "CANCELED"
)
