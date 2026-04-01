package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) Create(ctx context.Context, order model.OrderRequest) (model.OrderResponse, error) {
	if order.UserUUID == uuid.Nil {
		return model.OrderResponse{}, model.ErrBadRequest
	}

	if len(order.PartUUIDs) == 0 {
		return model.OrderResponse{}, model.ErrBadRequest
	}

	seen := make(map[uuid.UUID]struct{}, len(order.PartUUIDs))
	uuids := make([]string, 0, len(order.PartUUIDs))
	for _, v := range order.PartUUIDs {
		if v == uuid.Nil {
			return model.OrderResponse{}, model.ErrBadRequest
		}

		if _, ok := seen[v]; ok {
			return model.OrderResponse{}, model.ErrBadRequest
		}
		seen[v] = struct{}{}

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
