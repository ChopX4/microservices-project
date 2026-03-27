package part

import (
	"context"

	"github.com/ChopX4/raketka/inventory/internal/model"
	"github.com/ChopX4/raketka/inventory/internal/repository/converter"
)

func (r *repository) List(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage := make([]model.Part, 0)

	for _, v := range r.parts {
		p := converter.PartToModel(v)

		if len(filter.UUIDS) > 0 && !contains(filter.UUIDS, p.UUID) {
			continue
		}

		if len(filter.Names) > 0 && !contains(filter.Names, p.Name) {
			continue
		}

		if len(filter.Categories) > 0 && !contains(filter.Categories, p.Category) {
			continue
		}

		if len(filter.ManunufacturerCountries) > 0 && !contains(filter.ManunufacturerCountries, p.Manufacturer.Country) {
			continue
		}

		if len(filter.Tags) > 0 {
			match := false
			for _, ft := range filter.Tags {
				if contains(p.Tags, ft) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		storage = append(storage, p)
	}

	return storage, nil
}

func contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}

	return false
}
