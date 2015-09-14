package pezinventory

import "errors"

const (
	//InventoryCollectionName holds the name of the inventory collection
	InventoryCollectionName = "inventory"
	//LeaseCollectionName holds the name of the leases collection
	LeaseCollectionName = "leases"
	//InventoryItemStatusAvailable - this means the InventoryItem is in an available state
	InventoryItemStatusAvailable = "available"
	//InventoryItemStatusReserving - this means the InventoryItem is in a reserving state
	InventoryItemStatusReserving = "reserving"
	//InventoryItemStatusLeased = this means the InventoryItem is in a leased state
	InventoryItemStatusLeased = "leased"
)

var (
	//ErrInventoryNotAvailable indicates the given InventoryItem is either
	//not found or is not available to be leased.
	ErrInventoryNotAvailable = errors.New("The InventoryItem specified is not available")
)
