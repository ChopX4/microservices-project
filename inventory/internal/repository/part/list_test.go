package part

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ChopX4/raketka/inventory/internal/model"
	repoModel "github.com/ChopX4/raketka/inventory/internal/repository/model"
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

	ctx := context.Background()

	mongoDbContainer, err := mongodb.Run(ctx, "mongo:6.0")
	if err != nil {
		t.Fatalf("failed to terminate container: %v", err)
	}
	defer func() {
		if err := mongoDbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}()

	endpoint, err := mongoDbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("failed to disconnect: %v", err)
		}
	}()

	db := client.Database("test_inventory")

	repo, err := NewRepository(ctx, db)
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}

	initialData := []any{repoPartEngine, repoPartFuel, repoPartWing}
	_, err = db.Collection("parts").InsertMany(ctx, initialData)
	if err != nil {
		t.Fatalf("failed to seed data: %v", err)
	}

	tests := []struct {
		name      string
		filter    model.PartsFilter
		wantParts []model.Part
	}{
		{
			name: "Поиск по UUID (точное совпадение)",
			filter: model.PartsFilter{
				UUIDS: []string{"uuid-1", "uuid-3"},
			},
			wantParts: []model.Part{partEngine, partFuel},
		},
		{
			name: "Поиск по категориям",
			filter: model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine, model.CategoryWing},
			},
			wantParts: []model.Part{partEngine, partWing},
		},
		{
			name: "Поиск по стране",
			filter: model.PartsFilter{
				ManufacturerCountries: []string{"Germany"},
			},
			wantParts: []model.Part{partEngine},
		},
		{
			name: "Фильтр, который ничего не найдет",
			filter: model.PartsFilter{
				Names: []string{"Non-existent-name"},
			},
			wantParts: []model.Part{},
		},
		{
			name:      "Пустой фильтр - возвращаем все данные",
			filter:    model.PartsFilter{},
			wantParts: []model.Part{partEngine, partWing, partFuel},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ := repo.List(context.Background(), test.filter)

			assert.ElementsMatch(t, test.wantParts, got)
		})
	}
}
