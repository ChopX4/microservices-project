package order

import (
	"github.com/google/uuid"

	"github.com/ChopX4/raketka/order/internal/model"
)

func (s *service) validateCreateOrderRequest(order model.OrderRequest) ([]string, error) {
	if order.UserUUID == uuid.Nil {
		return nil, model.ErrBadRequest
	}

	if len(order.PartUUIDs) == 0 {
		return nil, model.ErrBadRequest
	}

	seen := make(map[uuid.UUID]struct{}, len(order.PartUUIDs))
	uuids := make([]string, 0, len(order.PartUUIDs))
	for _, partUUID := range order.PartUUIDs {
		if partUUID == uuid.Nil {
			return nil, model.ErrBadRequest
		}

		if _, ok := seen[partUUID]; ok {
			return nil, model.ErrBadRequest
		}

		seen[partUUID] = struct{}{}
		uuids = append(uuids, partUUID.String())
	}

	return uuids, nil
}

func (s *service) validateOrderUUID(orderUUID string) error {
	if !model.IsValidUUID(orderUUID) {
		return model.ErrBadRequest
	}

	return nil
}

func (s *service) validatePayOrderRequest(req model.PayOrderRequest) error {
	if !model.IsValidUUID(req.OrderUuid) {
		return model.ErrBadRequest
	}

	if !req.PaymentMethod.IsValid() {
		return model.ErrBadRequest
	}

	return nil
}

func (s *service) validateOrderStatusForPay(status model.OrderStatus) error {
	if status == model.OrderStatusCanceled ||
		status == model.OrderStatusPaid ||
		status == model.OrderStatusCompleted {
		return model.ErrConflict
	}

	return nil
}

func (s *service) validateOrderStatusForCancel(status model.OrderStatus) error {
	if status == model.OrderStatusCanceled ||
		status == model.OrderStatusPaid ||
		status == model.OrderStatusCompleted {
		return model.ErrConflict
	}

	return nil
}

func (s *service) validateOrderStatusForComplete(status model.OrderStatus) error {
	if status == model.OrderStatusCanceled ||
		status == model.OrderStatusCompleted ||
		status == model.OrderStatusPendingPayment {
		return model.ErrConflict
	}

	return nil
}
