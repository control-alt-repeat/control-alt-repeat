package internal

import (
	"context"

	aws "github.com/control-alt-repeat/control-alt-repeat/internal/aws"
	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func MoveItem(ctx context.Context, itemID string, newShelf string) error {
	var warehouseItem models.WarehouseItem

	warehouseItem, err := warehouse.GetWarehouseItem(ctx, itemID)
	if err != nil {
		return err
	}

	warehouseItem.Shelf = newShelf

	for _, ebayListingID := range warehouseItem.EbayListingIDs {
		err = ebay.ReviseSKU(ctx, ebayListingID, warehouseItem.ToEbaySKU())
		if err != nil {
			return err
		}

		err = aws.SaveJsonObjectS3(ctx, "control-alt-repeat-warehouse", warehouseItem.ControlAltRepeatID, warehouseItem)
		if err != nil {
			return err
		}
	}
	return nil
}
