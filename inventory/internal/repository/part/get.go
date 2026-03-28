package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ChopX4/raketka/inventory/internal/model"
	repoConverter "github.com/ChopX4/raketka/inventory/internal/repository/converter"
	repoModel "github.com/ChopX4/raketka/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	var part repoModel.Part

	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}

		return model.Part{}, err
	}

	return repoConverter.PartToModel(part), nil
}
