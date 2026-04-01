package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ChopX4/raketka/inventory/internal/converter"
	"github.com/ChopX4/raketka/inventory/internal/model"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

// ListParts обрабатывает запрос на получение списка деталей по фильтру.
func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	parts, err := a.inventoryService.List(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
		if errors.Is(err, model.ErrInvalidCategory) {
			return nil, status.Error(codes.InvalidArgument, "filter contains invalid category")
		}
		return nil, err
	}

	protoParts := make([]*inventory_v1.Part, 0, len(parts))
	for _, v := range parts {
		protoParts = append(protoParts, converter.PartToProto(v))
	}

	return &inventory_v1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
