package kafka

import "github.com/ChopX4/raketka/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
