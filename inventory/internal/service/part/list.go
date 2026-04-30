package part

import (
	"context"

	"github.com/ChopX4/raketka/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	if err := s.validateListFilter(filter); err != nil {
		return nil, err
	}

	parts, err := s.inventoryRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
