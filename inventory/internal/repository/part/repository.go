package part

import (
	"sync"

	"github.com/ChopX4/raketka/inventory/internal/repository/model"
)

type repository struct {
	mu    sync.RWMutex
	parts map[string]model.Part
}

func NewRepository() *repository {
	return &repository{
		parts: make(map[string]model.Part),
	}
}
