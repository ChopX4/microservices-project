package model

import "github.com/google/uuid"

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

func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodCard, PaymentMethodSPB, PaymentMethodCreditCard, PaymentMethodInvestorMoney:
		return true
	default:
		return false
	}
}

func IsValidUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
