package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/ChopX4/raketka/notification/internal/model"
	events_v1 "github.com/ChopX4/raketka/shared/pkg/proto/events/v1"
)

type orderDecoder struct{}

func NewOrderDecoder() *orderDecoder {
	return &orderDecoder{}
}

func (d *orderDecoder) Decode(data []byte) (model.OrderPaid, error) {
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
