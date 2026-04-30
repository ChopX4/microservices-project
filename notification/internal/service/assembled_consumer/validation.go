package assembledconsumer

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/ChopX4/raketka/notification/internal/model"
)

var errInvalidShipAssembledEvent = errors.New("invalid ship assembled event")

func (s *service) validateShipAssembledEvent(event model.ShipAssembled) error {
	if strings.TrimSpace(event.EventUuid) == "" {
		return errInvalidShipAssembledEvent
	}

	if strings.TrimSpace(event.OrderUuid) == "" {
		return errInvalidShipAssembledEvent
	}

	if strings.TrimSpace(event.UserUuid) == "" {
		return errInvalidShipAssembledEvent
	}

	if event.BuildTimeSec <= 0 {
		return errInvalidShipAssembledEvent
	}

	if _, err := uuid.Parse(event.EventUuid); err != nil {
		return errInvalidShipAssembledEvent
	}

	if _, err := uuid.Parse(event.OrderUuid); err != nil {
		return errInvalidShipAssembledEvent
	}

	if _, err := uuid.Parse(event.UserUuid); err != nil {
		return errInvalidShipAssembledEvent
	}

	return nil
}
