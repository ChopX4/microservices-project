package model

import "github.com/google/uuid"

type OrderByUUID struct {
	OrderUUID       uuid.UUID     `db:"order_uuid"`
	UserUUID        uuid.UUID     `db:"user_uuid"`
	PartUuids       []uuid.UUID   `db:"part_uuids"`
	TotalPrice      float32       `db:"total_price"`
	TransactionUUID uuid.UUID     `db:"transaction_uuid"`
	PaymentMethod   PaymentMethod `db:"payment_method"`
	Status          OrderStatus   `db:"status"`
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
