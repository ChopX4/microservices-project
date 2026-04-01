package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/ChopX4/raketka/inventory/internal/model"
	"github.com/ChopX4/raketka/inventory/internal/repository/converter"
	repoModel "github.com/ChopX4/raketka/inventory/internal/repository/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	mongoFilter := filterForList(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		logger.Error(ctx, "failed to list parts from mongo", zap.Any("filter", mongoFilter), zap.Error(err))
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			logger.Error(ctx, "failed to close cursor", zap.Error(cerr))
		}
	}()

	var parts []repoModel.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		logger.Error(ctx, "failed to decode parts from mongo cursor", zap.Any("filter", mongoFilter), zap.Error(err))
		return nil, err
	}

	storage := make([]model.Part, 0, len(parts))
	for _, v := range parts {
		modelPart := converter.PartToModel(v)
		storage = append(storage, modelPart)
	}

	return storage, nil
}

func filterForList(filter model.PartsFilter) bson.M {
	mongoFilter := bson.M{}

	if len(filter.UUIDS) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": filter.UUIDS}
	}

	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		mongoFilter["category"] = bson.M{"$in": filter.Categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	return mongoFilter
}
