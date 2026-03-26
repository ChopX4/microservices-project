package part

import (
	"testing"
	"time"

	"github.com/ChopX4/raketka/inventory/internal/model"
	repoModel "github.com/ChopX4/raketka/inventory/internal/repository/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestList(t *testing.T) {
	var (
		t1 = time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
		t2 = time.Date(2023, 2, 1, 11, 0, 0, 0, time.UTC)
	)

	// Деталь 1: Двигатель из Германии
	partEngine := model.Part{
		UUID:         "uuid-1",
		Name:         "V12-Turbo",
		Category:     model.CategoryEngine,
		Manufacturer: model.Manufacturer{Country: "Germany", Name: "BMW"},
		Tags:         []string{"power", "heavy"},
		CreatedAt:    &t1,
	}

	// Деталь 2: Крыло из Японии
	partWing := model.Part{
		UUID:         "uuid-2",
		Name:         "Carbon-Wing",
		Category:     model.CategoryWing,
		Manufacturer: model.Manufacturer{Country: "Japan", Name: "Mitsubishi"},
		Tags:         []string{"light", "aerodynamic"},
		CreatedAt:    &t2,
	}

	// Деталь 3: Топливный насос из Японии
	partFuel := model.Part{
		UUID:         "uuid-3",
		Name:         "Fuel-Pump-X",
		Category:     model.CategoryFuel,
		Manufacturer: model.Manufacturer{Country: "Japan", Name: "Denso"},
		Tags:         []string{"fuel", "power"},
		CreatedAt:    &t1,
	}

	// Деталь 1: Двигатель из Германии
	repoPartEngine := repoModel.Part{
		UUID:         "uuid-1",
		Name:         "V12-Turbo",
		Category:     repoModel.CategoryEngine,
		Manufacturer: repoModel.Manufacturer{Country: "Germany", Name: "BMW"},
		Tags:         []string{"power", "heavy"},
		CreatedAt:    &t1,
	}

	// Деталь 2: Крыло из Японии
	repoPartWing := repoModel.Part{
		UUID:         "uuid-2",
		Name:         "Carbon-Wing",
		Category:     repoModel.CategoryWing,
		Manufacturer: repoModel.Manufacturer{Country: "Japan", Name: "Mitsubishi"},
		Tags:         []string{"light", "aerodynamic"},
		CreatedAt:    &t2,
	}

	// Деталь 3: Топливный насос из Японии
	repoPartFuel := repoModel.Part{
		UUID:         "uuid-3",
		Name:         "Fuel-Pump-X",
		Category:     repoModel.CategoryFuel,
		Manufacturer: repoModel.Manufacturer{Country: "Japan", Name: "Denso"},
		Tags:         []string{"fuel", "power"},
		CreatedAt:    &t1,
	}

	tests := []struct {
		name        string
		filter      model.PartsFilter
		wantParts   []model.Part
		repoStorage map[string]repoModel.Part
	}{
		{
			name: "Поиск по UUID (точное совпадение)",
			filter: model.PartsFilter{
				UUIDS: []string{"uuid-1", "uuid-3"},
			},
			wantParts: []model.Part{partEngine, partFuel},
			repoStorage: map[string]repoModel.Part{
				repoPartEngine.UUID: repoPartEngine,
				repoPartWing.UUID:   repoPartWing,
				repoPartFuel.UUID:   repoPartFuel,
			},
		},
		{
			name: "Поиск по категориям",
			filter: model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine, model.CategoryWing},
			},
			wantParts: []model.Part{partEngine, partWing},
			repoStorage: map[string]repoModel.Part{
				repoPartEngine.UUID: repoPartEngine,
				repoPartWing.UUID:   repoPartWing,
				repoPartFuel.UUID:   repoPartFuel,
			},
		},
		{
			name: "Поиск по стране",
			filter: model.PartsFilter{
				ManunufacturerCountries: []string{"Germany"},
			},
			wantParts: []model.Part{partEngine},
			repoStorage: map[string]repoModel.Part{
				repoPartEngine.UUID: repoPartEngine,
				repoPartWing.UUID:   repoPartWing,
				repoPartFuel.UUID:   repoPartFuel,
			},
		},
		{
			name: "Фильтр, который ничего не найдет",
			filter: model.PartsFilter{
				Names: []string{"Non-existent-name"},
			},
			wantParts: []model.Part{},
			repoStorage: map[string]repoModel.Part{
				repoPartEngine.UUID: repoPartEngine,
				repoPartWing.UUID:   repoPartWing,
				repoPartFuel.UUID:   repoPartFuel,
			},
		},
		{
			name:      "Пустой фильтр - возвращаем все данные",
			filter:    model.PartsFilter{},
			wantParts: []model.Part{partEngine, partWing, partFuel},
			repoStorage: map[string]repoModel.Part{
				repoPartEngine.UUID: repoPartEngine,
				repoPartWing.UUID:   repoPartWing,
				repoPartFuel.UUID:   repoPartFuel,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repository{
				parts: test.repoStorage,
			}

			got, _ := repo.List(context.Background(), test.filter)

			assert.ElementsMatch(t, test.wantParts, got)
		})
	}
}
