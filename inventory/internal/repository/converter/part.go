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
		Metadata:      MetadataToModel(repo.Metadata),
		CreatedAt:     repo.CreatedAt,
		UpdatedAt:     repo.UpdatedAt,
	}
}

func CategoryToModel(repo repoModel.Category) model.Category {
	return model.Category(repo)
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

func MetadataToModel(repoMeta map[string]repoModel.Value) map[string]model.Value {
	if repoMeta == nil {
		return nil
	}

	result := make(map[string]model.Value)

	for key, val := range repoMeta {
		switch v := val.(type) {
		case repoModel.StringValue:
			result[key] = model.StringValue{V: v.V}
		case repoModel.Int64Value:
			result[key] = model.Int64Value{V: v.V}
		case repoModel.Float64Value:
			result[key] = model.Float64Value{V: v.V}
		case repoModel.BoolValue:
			result[key] = model.BoolValue{V: v.V}
		}
	}

	return result
}
