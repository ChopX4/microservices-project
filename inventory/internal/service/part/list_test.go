package part

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ChopX4/raketka/inventory/internal/model"
	"github.com/ChopX4/raketka/inventory/internal/repository/mocks"
)

func TestList(t *testing.T) {
	var (
		t1 = time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
		t2 = time.Date(2023, 2, 1, 11, 0, 0, 0, time.UTC)
	)

	partEngine := model.Part{
		UUID:         "550e8400-e29b-41d4-a716-446655440001",
		Name:         "V12-Turbo",
		Category:     model.CategoryEngine,
		Manufacturer: model.Manufacturer{Country: "Germany", Name: "BMW"},
		Tags:         []string{"power", "heavy"},
		CreatedAt:    &t1,
	}

	// Деталь 2: Крыло из Японии
	partWing := model.Part{
		UUID:         "550e8400-e29b-41d4-a716-446655440002",
		Name:         "Carbon-Wing",
		Category:     model.CategoryWing,
		Manufacturer: model.Manufacturer{Country: "Japan", Name: "Mitsubishi"},
		Tags:         []string{"light", "aerodynamic"},
		CreatedAt:    &t2,
	}

	// Деталь 3: Топливный насос из Японии
	partFuel := model.Part{
		UUID:         "550e8400-e29b-41d4-a716-446655440003",
		Name:         "Fuel-Pump-X",
		Category:     model.CategoryFuel,
		Manufacturer: model.Manufacturer{Country: "Japan", Name: "Denso"},
		Tags:         []string{"fuel", "power"},
		CreatedAt:    &t1,
	}

	tests := []struct {
		name        string
		filter      model.PartsFilter
		prepareMock func(ir *mocks.InventoryRepository, f model.PartsFilter)
		wantParts   []model.Part
	}{
		{
			name: "Поиск по UUID (точное совпадение)",
			filter: model.PartsFilter{
				UUIDS: []string{"550e8400-e29b-41d4-a716-446655440001", "550e8400-e29b-41d4-a716-446655440003"},
			},
			prepareMock: func(ir *mocks.InventoryRepository, f model.PartsFilter) {
				ir.On("List", context.Background(), f).Return([]model.Part{partEngine, partFuel}, nil)
			},
			wantParts: []model.Part{partEngine, partFuel},
		},
		{
			name: "Поиск по категориям",
			filter: model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine, model.CategoryWing},
			},
			prepareMock: func(ir *mocks.InventoryRepository, f model.PartsFilter) {
				ir.On("List", context.Background(), f).Return([]model.Part{partEngine, partWing}, nil)
			},
			wantParts: []model.Part{partEngine, partWing},
		},
		{
			name: "Поиск по стране",
			filter: model.PartsFilter{
				ManufacturerCountries: []string{"Germany"},
			},
			prepareMock: func(ir *mocks.InventoryRepository, f model.PartsFilter) {
				ir.On("List", context.Background(), f).Return([]model.Part{partEngine}, nil)
			},
			wantParts: []model.Part{partEngine},
		},
		{
			name: "Фильтр, который ничего не найдет",
			filter: model.PartsFilter{
				Names: []string{"Non-existent-name"},
			},
			prepareMock: func(ir *mocks.InventoryRepository, f model.PartsFilter) {
				ir.On("List", context.Background(), f).Return([]model.Part{}, nil)
			},
			wantParts: []model.Part{},
		},
		{
			name:   "Пустой фильтр - возвращаем все данные",
			filter: model.PartsFilter{},
			prepareMock: func(ir *mocks.InventoryRepository, f model.PartsFilter) {
				ir.On("List", context.Background(), f).Return([]model.Part{partEngine, partWing, partFuel}, nil)
			},
			wantParts: []model.Part{partEngine, partWing, partFuel},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.InventoryRepository)
			test.prepareMock(mockRepo, test.filter)

			s := service{
				inventoryRepository: mockRepo,
			}

			got, _ := s.List(context.Background(), test.filter)

			assert.ElementsMatch(t, test.wantParts, got)

			mockRepo.AssertExpectations(t)
		})
	}
}
