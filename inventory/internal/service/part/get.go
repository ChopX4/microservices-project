package part

import (
	"context"

	"github.com/ChopX4/raketka/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, partUUID string) (model.Part, error) {
	if err := s.validateGetRequest(partUUID); err != nil {
		return model.Part{}, err
	}

	part, err := s.inventoryRepository.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
