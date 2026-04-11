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

func TestListParts(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Millisecond)

	testParts := []model.Part{
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440000",
			Name:          "Rocket Engine",
			Description:   "Main engine",
			Price:         1200.5,
			StockQuantity: 3,
			Category:      model.CategoryEngine,
			CreatedAt:     &now,
		},
		{
			UUID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Fuel Pump",
			Description:   "Fuel system part",
			Price:         340.2,
			StockQuantity: 7,
			Category:      model.CategoryFuel,
			CreatedAt:     &now,
		},
	}

	expectedFilter := model.PartsFilter{
		UUIDS:      []string{"550e8400-e29b-41d4-a716-446655440000"},
		Categories: []model.Category{model.CategoryEngine},
	}

	tests := []struct {
		name       string
		req        *inventory_v1.ListPartsRequest
		prepare    func(mockService *serviceMocks.InventoryService)
		assertFunc func(t *testing.T, resp *inventory_v1.ListPartsResponse, err error)
	}{
		{
			name: "nil request",
			req:  nil,
			assertFunc: func(t *testing.T, resp *inventory_v1.ListPartsResponse, err error) {
				assert.Nil(t, resp)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Equal(t, "request is required", status.Convert(err).Message())
			},
		},
		{
			name: "success",
			req: &inventory_v1.ListPartsRequest{
				Filter: &inventory_v1.PartsFilter{
					Uuids:      []string{"550e8400-e29b-41d4-a716-446655440000"},
					Categories: []inventory_v1.Category{inventory_v1.Category_CATEGORY_ENGINE},
				},
			},
			prepare: func(mockService *serviceMocks.InventoryService) {
				mockService.On("List", ctx, expectedFilter).Return(testParts, nil)
			},
			assertFunc: func(t *testing.T, resp *inventory_v1.ListPartsResponse, err error) {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) && assert.Len(t, resp.Parts, 2) {
					assert.Equal(t, testParts[0].UUID, resp.Parts[0].GetUuid())
					assert.Equal(t, testParts[1].UUID, resp.Parts[1].GetUuid())
				}
			},
		},
		{
			name: "invalid category",
			req: &inventory_v1.ListPartsRequest{
				Filter: &inventory_v1.PartsFilter{},
			},
			prepare: func(mockService *serviceMocks.InventoryService) {
				mockService.On("List", ctx, model.PartsFilter{
					Categories: []model.Category{},
				}).Return(nil, model.ErrInvalidCategory)
			},
			assertFunc: func(t *testing.T, resp *inventory_v1.ListPartsResponse, err error) {
				assert.Nil(t, resp)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				assert.Equal(t, "filter contains invalid category", status.Convert(err).Message())
			},
		},
		{
			name: "unexpected service error",
			req: &inventory_v1.ListPartsRequest{
				Filter: &inventory_v1.PartsFilter{},
			},
			prepare: func(mockService *serviceMocks.InventoryService) {
				mockService.On("List", ctx, model.PartsFilter{
					Categories: []model.Category{},
				}).Return(nil, assert.AnError)
			},
			assertFunc: func(t *testing.T, resp *inventory_v1.ListPartsResponse, err error) {
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
			resp, err := apiHandler.ListParts(ctx, test.req)

			test.assertFunc(t, resp, err)
		})
	}
}
