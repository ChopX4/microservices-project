package order

import (
	"context"
	"testing"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/order/internal/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCancel(t *testing.T) {
	var (
		orderID = uuid.New()
		userID  = uuid.New()
		part1   = uuid.New()
		part2   = uuid.New()
		transID = uuid.New()
	)

	testOrder := model.OrderByUUID{
		OrderUUID:       orderID,
		UserUUID:        userID,
		PartUuids:       []uuid.UUID{part1, part2},
		TotalPrice:      1500.50,
		TransactionUUID: transID,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.OrderStatusPaid,
	}

	testOrderConflict := model.OrderByUUID{
		OrderUUID:       orderID,
		UserUUID:        userID,
		PartUuids:       []uuid.UUID{part1, part2},
		TotalPrice:      1500.50,
		TransactionUUID: transID,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.OrderStatusCanceled,
	}

	testOrderCanceled := model.OrderByUUID{
		OrderUUID:       orderID,
		UserUUID:        userID,
		PartUuids:       []uuid.UUID{part1, part2},
		TotalPrice:      1500.50,
		TransactionUUID: transID,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.OrderStatusCanceled,
	}

	tests := []struct {
		name        string
		orderUuid   string
		prepareMock func(or *mocks.OrderRepository, uuid string)
		expError    error
	}{
		{
			name:      "Успешная отмена",
			orderUuid: orderID.String(),
			prepareMock: func(or *mocks.OrderRepository, uuid string) {
				or.On("Get", uuid).Return(testOrder, nil)
				or.On("Update", context.Background(), testOrderCanceled).Return(nil)
			},
			expError: nil,
		},
		{
			name:      "Заказ не найден",
			orderUuid: "testid",
			prepareMock: func(or *mocks.OrderRepository, uuid string) {
				or.On("Get", uuid).Return(model.OrderByUUID{}, model.ErrNotFound)
			},
			expError: model.ErrNotFound,
		},
		{
			name:      "Ошибка конфликт",
			orderUuid: orderID.String(),
			prepareMock: func(or *mocks.OrderRepository, uuid string) {
				or.On("Get", uuid).Return(testOrderConflict, nil)
			},
			expError: model.ErrConflict,
		},
		{
			name:      "Ошибка при апдейте",
			orderUuid: orderID.String(),
			prepareMock: func(or *mocks.OrderRepository, uuid string) {
				or.On("Get", uuid).Return(testOrder, nil)
				or.On("Update", context.Background(), testOrderCanceled).Return(model.ErrNotFound)
			},
			expError: model.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.OrderRepository)
			test.prepareMock(mockRepo, test.orderUuid)

			s := service{
				orderRepository: mockRepo,
			}

			err := s.Cancel(context.Background(), test.orderUuid)
			if err != nil {
				assert.ErrorIs(t, err, test.expError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
