package order

import (
	"sync"

	"github.com/ChopX4/raketka/order/internal/repository/model"
)

type repository struct {
	mu     sync.RWMutex
	orders map[string]model.OrderByUUID
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]model.OrderByUUID),
	}
}
