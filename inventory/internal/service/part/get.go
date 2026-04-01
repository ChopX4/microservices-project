package part

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, partUUID string) (model.Part, error) {
	if strings.TrimSpace(partUUID) == "" {
		return model.Part{}, model.ErrInvalidUUID
	}

	if _, err := uuid.Parse(partUUID); err != nil {
		return model.Part{}, model.ErrInvalidUUID
	}

	part, err := s.inventoryRepository.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
