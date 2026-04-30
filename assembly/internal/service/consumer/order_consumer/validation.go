package orderconsumer

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/assembly/internal/model"
)

var errInvalidOrderPaidEvent = errors.New("invalid order paid event")

func (s *service) validateOrderPaidEvent(event model.OrderPaid) error {
	if strings.TrimSpace(event.EventUuid) == "" {
		return errInvalidOrderPaidEvent
	}

	if strings.TrimSpace(event.OrderUuid) == "" {
		return errInvalidOrderPaidEvent
	}

	if strings.TrimSpace(event.UserUuid) == "" {
		return errInvalidOrderPaidEvent
	}

	if strings.TrimSpace(event.TransactionUuid) == "" {
		return errInvalidOrderPaidEvent
	}

	if _, err := uuid.Parse(event.EventUuid); err != nil {
		return errInvalidOrderPaidEvent
	}

	if _, err := uuid.Parse(event.OrderUuid); err != nil {
		return errInvalidOrderPaidEvent
	}

	if _, err := uuid.Parse(event.UserUuid); err != nil {
		return errInvalidOrderPaidEvent
	}

	if _, err := uuid.Parse(event.TransactionUuid); err != nil {
		return errInvalidOrderPaidEvent
	}

	return nil
}
