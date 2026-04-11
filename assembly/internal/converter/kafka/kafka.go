package kafka

import "github.com/ChopX4/raketka/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}
