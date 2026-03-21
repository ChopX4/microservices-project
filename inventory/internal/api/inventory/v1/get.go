package v1

import (
	"context"
	"errors"

	"github.com/ChopX4/raketka/inventory/internal/converter"
	"github.com/ChopX4/raketka/inventory/internal/model"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	part, err := a.inventoryService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &inventory_v1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
