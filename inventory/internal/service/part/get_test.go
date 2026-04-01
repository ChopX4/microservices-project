package part

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ChopX4/raketka/inventory/internal/model"
	"github.com/ChopX4/raketka/inventory/internal/repository/mocks"
)

func TestGet(t *testing.T) {
	mockRepo := new(mocks.InventoryRepository)
	s := service{
		inventoryRepository: mockRepo,
	}

	now := time.Now()
	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	testPart := model.Part{
		UUID:          testUUID,
		Name:          "Rocket",
		Description:   "Rocket12",
		Price:         12322.12,
		StockQuantity: 10,
		CreatedAt:     &now,
	}

	tests := []struct {
		name     string
		testUuid string
		preMock  func()
		expPart  model.Part
		expError error
	}{
		{
			name:     "детать найдена",
			testUuid: testUUID,
			preMock: func() {
				mockRepo.On("Get", context.Background(), testUUID).Return(testPart, nil)
			},
			expPart:  testPart,
			expError: nil,
		},
		{
			name:     "Деталь не найдена",
			testUuid: "550e8400-e29b-41d4-a716-446655440001",
			preMock: func() {
				mockRepo.On("Get", context.Background(), "550e8400-e29b-41d4-a716-446655440001").Return(model.Part{}, model.ErrPartNotFound)
			},
			expPart:  model.Part{},
			expError: model.ErrPartNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.preMock()

			got, err := s.Get(context.Background(), test.testUuid)
			if err != nil {
				assert.ErrorIs(t, err, test.expError)
				assert.Equal(t, test.expPart, got)
			} else {
				assert.Equal(t, test.expPart, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
