package decoder

import (
	"github.com/ChopX4/raketka/assembly/internal/model"
	events_v1 "github.com/ChopX4/raketka/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaid, error) {
	var result events_v1.OrderPaid
	if err := proto.Unmarshal(data, &result); err != nil {
		return model.OrderPaid{}, err
	}

	return model.OrderPaid{
		EventUuid:       result.GetEventUuid(),
		OrderUuid:       result.GetOrderUuid(),
		UserUuid:        result.GetUserUuid(),
		TransactionUuid: result.GetTransactionUuid(),
	}, nil
}
