package v1

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ChopX4/raketka/inventory/internal/model"
	serviceMocks "github.com/ChopX4/raketka/inventory/internal/service/mocks"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

func TestGetPart(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)
	partUUID := "550e8400-e29b-41d4-a716-446655440000"

	testPart := model.Part{
		UUID:          partUUID,
		Name:          "Rocket Engine",
		Description:   "Main engine",
		Price:         1200.5,
		StockQuantity: 3,
		Category:      model.CategoryEngine,
		CreatedAt:     &now,
	}

	tests := []struct {
		name       string
		req        *inventory_v1.GetPartRequest
		prepare    func(mockService *serviceMocks.InventoryService)
		assertFunc func(t *testing.T, resp *inventory_v1.GetPartResponse, err error)
	}{
		{
			name: "nil request",
			req:  nil,
			assertFunc: func(t *testing.T, resp *inventory_v1.GetPartResponse, err error) {
				assert.Nil(t, resp)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Equal(t, "request is required", status.Convert(err).Message())
			},
		},
		{
			name: "success",
			req:  &inventory_v1.GetPartRequest{Uuid: partUUID},
			prepare: func(mockService *serviceMocks.InventoryService) {
				mockService.On("Get", ctx, partUUID).Return(testPart, nil)
			},
			assertFunc: func(t *testing.T, resp *inventory_v1.GetPartResponse, err error) {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) && assert.NotNil(t, resp.Part) {
					assert.Equal(t, testPart.UUID, resp.Part.GetUuid())
					assert.Equal(t, testPart.Name, resp.Part.GetName())
					assert.Equal(t, testPart.Description, resp.Part.GetDescription())
					assert.Equal(t, testPart.Price, resp.Part.GetPrice())
					assert.Equal(t, testPart.StockQuantity, resp.Part.GetStockQuantity())
				}
			},
		},
		{
			name: "invalid uuid",
			req:  &inventory_v1.GetPartRequest{Uuid: "bad-uuid"},
			prepare: func(mockService *serviceMocks.InventoryService) {
				mockService.On("Get", ctx, "bad-uuid").Return(model.Part{}, model.ErrInvalidUUID)
			},
			assertFunc: func(t *testing.T, resp *inventory_v1.GetPartResponse, err error) {
				assert.Nil(t, resp)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Equal(t, "uuid must be a valid UUID", status.Convert(err).Message())
			},
		},
		{
			name: "part not found",
			req:  &inventory_v1.GetPartRequest{Uuid: partUUID},
			prepare: func(mockService *serviceMocks.InventoryService) {
				mockService.On("Get", ctx, partUUID).Return(model.Part{}, model.ErrPartNotFound)
			},
			assertFunc: func(t *testing.T, resp *inventory_v1.GetPartResponse, err error) {
				assert.Nil(t, resp)
				assert.Equal(t, codes.NotFound, status.Code(err))
				assert.Equal(t, "part with UUID "+partUUID+" not found", status.Convert(err).Message())
			},
		},
		{
			name: "unexpected service error",
			req:  &inventory_v1.GetPartRequest{Uuid: partUUID},
			prepare: func(mockService *serviceMocks.InventoryService) {
				mockService.On("Get", ctx, partUUID).Return(model.Part{}, assert.AnError)
			},
			assertFunc: func(t *testing.T, resp *inventory_v1.GetPartResponse, err error) {
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService := serviceMocks.NewInventoryService(t)
			if test.prepare != nil {
				test.prepare(mockService)
			}

			apiHandler := NewApi(mockService)
			resp, err := apiHandler.GetPart(ctx, test.req)

			test.assertFunc(t, resp, err)
		})
	}
}
