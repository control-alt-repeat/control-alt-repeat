package internal

import (
	"context"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/warehouse"
)

func LookupItem(ctx context.Context, itemID string) (warehouse.WarehouseItem, []warehouse.EbayItemInternal, error) {
	warehouseItem, err := warehouse.GetWarehouseItem(ctx, itemID)
	if err != nil {
		return warehouse.WarehouseItem{}, nil, err
	}

	ebayInternalItems, err := warehouse.GetEbayInternalItems(ctx, warehouseItem.EbayListingIDs)
	if err != nil {
		return warehouse.WarehouseItem{}, nil, err
	}

	return warehouseItem, ebayInternalItems, err
}
