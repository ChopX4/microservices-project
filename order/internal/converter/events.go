package converter

import (
	"github.com/ChopX4/raketka/order/internal/model"
	events_v1 "github.com/ChopX4/raketka/shared/pkg/proto/events/v1"
)

func OrderPaidToProto(event model.OrderPaid) *events_v1.OrderPaid {
	return &events_v1.OrderPaid{
		EventUuid:       event.EventUuid,
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		TransactionUuid: event.TransactionUuid,
	}
}

func ShipAssembledToProto(event model.ShipAssembled) *events_v1.ShipAssembled {
	return &events_v1.ShipAssembled{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}
}
