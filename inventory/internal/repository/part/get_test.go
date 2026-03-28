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

func TestGet(t *testing.T) {
	uuid := "129318djasd123jasd123"
	timeNow := time.Now().UTC().Truncate(time.Millisecond)

	testPart := model.Part{
		UUID:          uuid,
		Name:          "Rocket",
		Description:   "Rocket12",
		Price:         12322.12,
		StockQuantity: 10,
		CreatedAt:     &timeNow,
	}

	repoPart := repoModel.Part{
		UUID:          uuid,
		Name:          "Rocket",
		Description:   "Rocket12",
		Price:         12322.12,
		StockQuantity: 10,
		CreatedAt:     &timeNow,
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

	repo, err := NewRepository(db)
	if err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}

	initialData := []any{repoPart}
	_, err = db.Collection("parts").InsertMany(ctx, initialData)
	if err != nil {
		t.Fatalf("failed to seed data: %v", err)
	}

	tests := []struct {
		name      string
		uuid      string
		wantPart  model.Part
		wantError error
	}{
		{
			name:      "Деталь найдена",
			uuid:      uuid,
			wantPart:  testPart,
			wantError: nil,
		},
		{
			name:      "Ошибка: Деталь не найдена",
			uuid:      "unknown",
			wantPart:  model.Part{},
			wantError: model.ErrPartNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.Get(context.Background(), test.uuid)

			if err != nil {
				assert.ErrorIs(t, err, test.wantError)
				assert.Equal(t, test.wantPart, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wantPart, got)
			}
		})
	}
}
