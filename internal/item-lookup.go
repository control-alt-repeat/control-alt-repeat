package internal

import (
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/warehouse"
)

func LookupItem(itemID string) (warehouse.WarehouseItem, []warehouse.EbayItemInternal, error) {
	warehouseItem, err := warehouse.GetWarehouseItem(itemID)
	if err != nil {
		return warehouse.WarehouseItem{}, nil, err
	}

	ebayInternalItems, err := warehouse.GetEbayInternalItems(warehouseItem.EbayListingIDs)
	if err != nil {
		return warehouse.WarehouseItem{}, nil, err
	}

	return warehouseItem, ebayInternalItems, err
}
