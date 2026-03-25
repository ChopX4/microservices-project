package converter

import (
	"github.com/ChopX4/raketka/order/internal/model"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

func PartsFilterToModel(filter *inventory_v1.PartsFilter) model.PartsFilter {
	categoriesStorage := make([]model.Category, 0, len(filter.Categories))

	protoCats := filter.GetCategories()
	for _, v := range protoCats {
		categoriesStorage = append(categoriesStorage, CategoryToModel(v))
	}

	return model.PartsFilter{
		UUIDS:                   filter.GetUuids(),
		Names:                   filter.GetNames(),
		Categories:              categoriesStorage,
		ManunufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                    filter.GetTags(),
	}
}

func CategoryToModel(category inventory_v1.Category) model.Category {
	return model.Category(category)
}
