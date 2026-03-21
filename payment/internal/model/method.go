package model

type PaymentMethod int32

const (
	PaymentMethodUnknown       PaymentMethod = 0
	PaymentMethodCard          PaymentMethod = 1
	PaymentMethodSPB           PaymentMethod = 2
	PaymentMethodCreditCard    PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)

type PayOrderRequest struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}
