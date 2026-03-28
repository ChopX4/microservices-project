package part

import (
	"context"

	"github.com/ChopX4/raketka/inventory/internal/model"
	repoConverter "github.com/ChopX4/raketka/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(part), nil
}
