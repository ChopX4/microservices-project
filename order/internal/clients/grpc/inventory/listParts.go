package inventory

import (
	"context"

	"go.uber.org/zap"

	"github.com/ChopX4/raketka/order/internal/clients/converter"
	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/ChopX4/raketka/platform/pkg/logger"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

func (c *inventoryClient) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	grpcReq := &inventory_v1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	}

	grpsReq, err := c.generatedClient.ListParts(ctx, grpcReq)
	if err != nil {
		logger.Error(ctx, "failed to list parts via inventory grpc", zap.Error(err))
		return nil, err
	}

	inventoryParts := grpsReq.GetParts()

	modelParts := make([]model.Part, 0, len(inventoryParts))

	for _, v := range inventoryParts {
		modelParts = append(modelParts, converter.PartToModel(v))
	}

	return modelParts, nil
}
