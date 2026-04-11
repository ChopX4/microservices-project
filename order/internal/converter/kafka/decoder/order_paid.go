package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/ChopX4/raketka/order/internal/model"
	events_v1 "github.com/ChopX4/raketka/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewShipDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembled, error) {
	var result events_v1.ShipAssembled
	if err := proto.Unmarshal(data, &result); err != nil {
		return model.ShipAssembled{}, err
	}

	return model.ShipAssembled{
		EventUuid:    result.GetEventUuid(),
		OrderUuid:    result.GetOrderUuid(),
		UserUuid:     result.GetUserUuid(),
		BuildTimeSec: result.GetBuildTimeSec(),
	}, nil
}
