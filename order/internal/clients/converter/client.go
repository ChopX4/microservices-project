package converter

import (
	"time"

	"github.com/ChopX4/raketka/order/internal/model"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ChopX4/raketka/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToProto(req model.PayOrderRequest) *payment_v1.PayOrderRequest {
	return &payment_v1.PayOrderRequest{
		OrderUuid:     req.OrderUuid,
		UserUuid:      req.UserUuid,
		PaymentMethod: PaymentMethodToProto(req.PaymentMethod),
	}
}

func PartsFilterToModel(filter *inventory_v1.PartsFilter) model.PartsFilter {
	categoriesStorage := make([]model.Category, 0, len(filter.GetCategories()))

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

func PartsFilterToProto(filter model.PartsFilter) *inventory_v1.PartsFilter {
	categoriesStorage := make([]inventory_v1.Category, 0, len(filter.Categories))

	modelCats := filter.Categories
	for _, v := range modelCats {
		categoriesStorage = append(categoriesStorage, CategoryToProto(v))
	}

	return &inventory_v1.PartsFilter{
		Uuids:                 filter.UUIDS,
		Names:                 filter.Names,
		Categories:            categoriesStorage,
		ManufacturerCountries: filter.ManunufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func PartToModel(part *inventory_v1.Part) model.Part {
	var createdAt *time.Time
	if part.CreatedAt != nil {
		tmp := part.CreatedAt.AsTime()
		createdAt = &tmp
	}

	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		tmp := part.UpdatedAt.AsTime()
		updatedAt = &tmp
	}

	return model.Part{
		UUID:          part.GetUuid(),
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      CategoryToModel(part.GetCategory()),
		Dimensions:    DimensionsToModel(part.GetDimensions()),
		Manufacturer:  ManufacturerToModel(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      MetadataToModel(part.GetMetadata()),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func CategoryToModel(category inventory_v1.Category) model.Category {
	return model.Category(category)
}

func CategoryToProto(category model.Category) inventory_v1.Category {
	return inventory_v1.Category(category)
}

func DimensionsToModel(dimensions *inventory_v1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.GetLength(),
		Width:  dimensions.GetWidth(),
		Height: dimensions.GetHeight(),
		Weight: dimensions.GetWeight(),
	}
}

func ManufacturerToModel(manufacturer *inventory_v1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.GetName(),
		Country: manufacturer.GetCountry(),
		WebSite: manufacturer.GetWebsite(),
	}
}

func MetadataToModel(meta map[string]*inventory_v1.Value) map[string]model.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]model.Value)

	for key, val := range meta {
		if val == nil || val.Kind == nil {
			continue
		}

		switch v := val.Kind.(type) {
		case *inventory_v1.Value_StringValue:
			result[key] = model.StringValue{V: v.StringValue}

		case *inventory_v1.Value_Int64Value:
			result[key] = model.Int64Value{V: v.Int64Value}

		case *inventory_v1.Value_DoubleValue:
			result[key] = model.Float64Value{V: v.DoubleValue}

		case *inventory_v1.Value_BoolValue:
			result[key] = model.BoolValue{V: v.BoolValue}
		}
	}

	return result
}

func PaymentMethodToProto(method model.PaymentMethod) payment_v1.PaymentMethod {
	switch method {
	case model.PaymentMethodUnknown:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSPB:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SPB
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	}

	return payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
}
