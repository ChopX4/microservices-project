package kafka

import "github.com/ChopX4/raketka/notification/internal/model"

type OrderDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}

type AssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
