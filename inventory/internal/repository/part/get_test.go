package part

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ChopX4/raketka/inventory/internal/model"
	repoModel "github.com/ChopX4/raketka/inventory/internal/repository/model"
)

func TestGet(t *testing.T) {
	uuid := "129318djasd123jasd123"
	timeNow := time.Now()

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

	tests := []struct {
		name      string
		repoData  map[string]repoModel.Part
		uuid      string
		wantPart  model.Part
		wantError error
	}{
		{
			name: "Деталь найдена",
			repoData: map[string]repoModel.Part{
				uuid: repoPart,
			},
			uuid:      uuid,
			wantPart:  testPart,
			wantError: nil,
		},
		{
			name:      "Ошибка: Деталь не найдена",
			repoData:  make(map[string]repoModel.Part),
			uuid:      "unknown",
			wantPart:  model.Part{},
			wantError: model.ErrPartNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repository{
				parts: test.repoData,
			}

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
