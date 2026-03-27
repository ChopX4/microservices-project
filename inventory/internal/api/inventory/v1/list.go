package v1

import (
	"context"

	"github.com/ChopX4/raketka/inventory/internal/converter"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts, err := a.inventoryService.List(ctx, converter.PartsFilterToModel(req.GetFilter()))
	if err != nil {
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
