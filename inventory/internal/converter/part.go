package converter

import (
	"github.com/ChopX4/raketka/inventory/internal/model"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartToProto(part model.Part) *inventory_v1.Part {
	var createdAt *timestamppb.Timestamp
	if part.CreatedAt != nil {
		createdAt = timestamppb.New(*part.CreatedAt)
	}
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventory_v1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToProto(part.Category),
		Dimensions:    DimensionsToProto(part.Dimensions),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToProto(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

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

func CategoryToProto(category model.Category) inventory_v1.Category {
	return inventory_v1.Category(category)
}

func CategoryToModel(category inventory_v1.Category) model.Category {
	return model.Category(category)
}

func DimensionsToProto(dimensions model.Dimensions) *inventory_v1.Dimensions {
	return &inventory_v1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToProto(manufacturer model.Manufacturer) *inventory_v1.Manufacturer {
	return &inventory_v1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.WebSite,
	}
}

func MetadataToProto(meta map[string]model.Value) map[string]*inventory_v1.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]*inventory_v1.Value)

	for key, val := range meta {
		protoVal := &inventory_v1.Value{}

		switch v := val.(type) {
		case model.StringValue:
			// В твоем коде: структура Value_StringValue, поле StringValue
			protoVal.Kind = &inventory_v1.Value_StringValue{StringValue: v.V}
		case model.Int64Value:
			// В твоем коде: структура Value_Int64Value, поле Int64Value
			protoVal.Kind = &inventory_v1.Value_Int64Value{Int64Value: v.V}
		case model.Float64Value:
			// В твоем коде: структура Value_DoubleValue (у тебя в прото Double, а не Float64)
			protoVal.Kind = &inventory_v1.Value_DoubleValue{DoubleValue: v.V}
		case model.BoolValue:
			// В твоем коде: структура Value_BoolValue, поле BoolValue
			protoVal.Kind = &inventory_v1.Value_BoolValue{BoolValue: v.V}
		}

		result[key] = protoVal
	}

	return result
}
