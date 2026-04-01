package converter

import (
	"github.com/ChopX4/raketka/inventory/internal/model"
	repoModel "github.com/ChopX4/raketka/inventory/internal/repository/model"
)

func PartToModel(repo repoModel.Part) model.Part {
	return model.Part{
		UUID:          repo.UUID,
		Name:          repo.Name,
		Description:   repo.Description,
		Price:         repo.Price,
		StockQuantity: repo.StockQuantity,
		Category:      CategoryToModel(repo.Category),
		Dimensions:    DimensionsToModel(repo.Dimensions),
		Manufacturer:  ManufacturerToModel(repo.Manufacturer),
		Tags:          repo.Tags,
		Metadata:      repo.Metadata,
		CreatedAt:     repo.CreatedAt,
		UpdatedAt:     repo.UpdatedAt,
	}
}

func CategoryToModel(repo repoModel.Category) model.Category {
	category := model.Category(repo)
	if !category.IsValid() {
		return model.CategoryUnknown
	}

	return category
}

func DimensionsToModel(repo repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: repo.Length,
		Width:  repo.Width,
		Height: repo.Height,
		Weight: repo.Weight,
	}
}

func ManufacturerToModel(repo repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    repo.Name,
		Country: repo.Country,
		WebSite: repo.WebSite,
	}
}
