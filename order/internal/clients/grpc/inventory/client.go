package inventory

import inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"

type inventoryClient struct {
	generatedClient inventory_v1.InventoryServiceClient
}

func NewInventoryClient(generatedClient inventory_v1.InventoryServiceClient) *inventoryClient {
	return &inventoryClient{
		generatedClient: generatedClient,
	}
}
