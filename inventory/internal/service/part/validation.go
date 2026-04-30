package part

import (
	"strings"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/inventory/internal/model"
)

func (s *service) validateGetRequest(partUUID string) error {
	if strings.TrimSpace(partUUID) == "" {
		return model.ErrInvalidUUID
	}

	if _, err := uuid.Parse(partUUID); err != nil {
		return model.ErrInvalidUUID
	}

	return nil
}

func (s *service) validateListFilter(filter model.PartsFilter) error {
	for _, partUUID := range filter.UUIDS {
		if strings.TrimSpace(partUUID) == "" {
			return model.ErrInvalidUUID
		}

		if _, err := uuid.Parse(partUUID); err != nil {
			return model.ErrInvalidUUID
		}
	}

	for _, category := range filter.Categories {
		if !category.IsValid() {
			return model.ErrInvalidCategory
		}
	}

	return nil
}
