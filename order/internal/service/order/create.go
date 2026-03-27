package order

import (
	"context"

	"github.com/ChopX4/raketka/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, order model.OrderRequest) (model.OrderResponse, error) {
	uuids := make([]string, 0, len(order.PartUUIDs))
	for _, v := range order.PartUUIDs {
		uuids = append(uuids, v.String())
	}

	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{UUIDS: uuids})
	if err != nil {
		return model.OrderResponse{}, err
	}

	if len(parts) != len(order.PartUUIDs) {
		return model.OrderResponse{}, model.ErrBadRequest
	}

	var totalPrice float64
	orderUUID := uuid.New()

	for _, v := range parts {
		totalPrice += v.Price
	}

	repoOrder := model.OrderByUUID{
		OrderUUID:  orderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: float32(totalPrice),
		Status:     model.OrderStatusPendingPayment,
	}
	if err := s.orderRepository.Create(ctx, repoOrder); err != nil {
		return model.OrderResponse{}, err
	}

	return model.OrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: float32(totalPrice),
	}, nil
}
