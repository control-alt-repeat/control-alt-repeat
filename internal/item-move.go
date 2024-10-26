package internal

import (
	"context"
	"fmt"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/warehouse"
)

func MoveItem(ctx context.Context, itemID string, newShelf string) error {
	var warehouseItem warehouse.WarehouseItem

	err := aws.LoadJsonObjectS3(ctx, warehouse.WarehouseItemsBucketName, itemID, &warehouseItem)

	if err != nil {
		return err
	}

	fmt.Printf("Loaded item '%s' from warehouse\n", warehouseItem.ControlAltRepeatID)

	warehouseItem.Shelf = newShelf

	for _, ebayListingID := range warehouseItem.EbayListingIDs {
		err = ebay.ReviseSKU(ctx, ebayListingID, warehouseItem.ToEbaySKU())
		if err != nil {
			return err
		}

		err = aws.SaveJsonObjectS3(ctx, warehouse.WarehouseItemsBucketName, warehouseItem.ControlAltRepeatID, warehouseItem)
		if err != nil {
			return err
		}
	}
	return nil
}
