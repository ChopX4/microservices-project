package v1

import (
	"github.com/ChopX4/raketka/inventory/internal/service"
	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewApi(inventoryService service.InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
