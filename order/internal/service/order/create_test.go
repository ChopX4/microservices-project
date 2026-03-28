package order

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ChopX4/raketka/order/internal/clients/grpc/mocks"
	"github.com/ChopX4/raketka/order/internal/model"
	repoMocks "github.com/ChopX4/raketka/order/internal/repository/mocks"
)

func TestCreate(t *testing.T) {
	var (
		id1      = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		id2      = uuid.MustParse("00000000-0000-0000-0000-000000000002")
		UserUUID = uuid.New()
		Parts    = []uuid.UUID{id1, id2}
	)

	errInventoryClient := errors.New("Какая-то ошибка из инвентори апи")
	errOrderCreate := errors.New("Ошибка при создании заказа в репо")

	testReq := model.OrderRequest{
		UserUUID:  UserUUID,
		PartUUIDs: Parts,
	}

	// Данные для ответа инвентаря
	partEngine := model.Part{UUID: id1.String(), Price: 200.00}
	partWing := model.Part{UUID: id2.String(), Price: 100.00}

	tests := []struct {
		name              string
		orderReq          model.OrderRequest
		prepareClientMock func(ic *mocks.InventoryClient)
		prepareRepoMock   func(or *repoMocks.OrderRepository)
		expTotalPrice     float64
		expError          error
	}{
		{
			name:     "Заказ успешно создан",
			orderReq: testReq,
			prepareClientMock: func(ic *mocks.InventoryClient) {
				ic.On("ListParts", mock.Anything, mock.Anything).
					Return([]model.Part{partEngine, partWing}, nil).Once()
			},
			prepareRepoMock: func(or *repoMocks.OrderRepository) {
				or.On("Create", mock.Anything, mock.MatchedBy(func(o model.OrderByUUID) bool {
					return o.UserUUID == UserUUID && o.TotalPrice == 300.00 && o.Status == model.OrderStatusPendingPayment
				})).Return(nil).Once()
			},
			expTotalPrice: 300.0,
			expError:      nil,
		},
		{
			name:     "Ошибка на inventory client",
			orderReq: testReq,
			prepareClientMock: func(ic *mocks.InventoryClient) {
				ic.On("ListParts", mock.Anything, mock.Anything).
					Return(nil, errInventoryClient).Once()
			},
			prepareRepoMock: func(or *repoMocks.OrderRepository) {},
			expError:        errInventoryClient,
		},
		{
			name:     "Ошибка неправильный ввод",
			orderReq: testReq,
			prepareClientMock: func(ic *mocks.InventoryClient) {
				ic.On("ListParts", mock.Anything, mock.Anything).
					Return([]model.Part{partEngine, partWing, {}}, nil).Once()
			},
			prepareRepoMock: func(or *repoMocks.OrderRepository) {},
			expError:        model.ErrBadRequest,
		},
		{
			name:     "Ошибка при создании на репо слое",
			orderReq: testReq,
			prepareClientMock: func(ic *mocks.InventoryClient) {
				ic.On("ListParts", mock.Anything, mock.Anything).
					Return([]model.Part{partEngine, partWing}, nil).Once()
			},
			prepareRepoMock: func(or *repoMocks.OrderRepository) {
				or.On("Create", mock.Anything, mock.MatchedBy(func(o model.OrderByUUID) bool {
					return o.UserUUID == UserUUID && o.TotalPrice == 300.00
				})).Return(errOrderCreate).Once()
			},
			expError: errOrderCreate,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(mocks.InventoryClient)
			test.prepareClientMock(mockClient)

			mockRepo := new(repoMocks.OrderRepository)
			test.prepareRepoMock(mockRepo)

			s := service{
				orderRepository: mockRepo,
				inventoryClient: mockClient,
			}

			got, err := s.Create(context.Background(), test.orderReq)

			if test.expError != nil {
				assert.ErrorIs(t, err, test.expError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, float32(test.expTotalPrice), got.TotalPrice)
				assert.NotEqual(t, uuid.Nil, got.OrderUUID)
			}

			mockClient.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
