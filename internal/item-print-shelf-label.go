package internal

import (
	"context"
	"fmt"
	"strings"

	aws "github.com/control-alt-repeat/control-alt-repeat/internal/aws"
	"github.com/control-alt-repeat/control-alt-repeat/internal/labels"
	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func ItemPrintShelfLabel(ctx context.Context, itemID string) error {
	var warehouseItem models.WarehouseItem
	var ebayItemInternal warehouse.EbayItemInternal

	warehouseItem, err := warehouse.GetWarehouseItem(ctx, itemID)
	if err != nil {
		return err
	}

	err = aws.LoadJsonObjectS3(ctx, warehouse.EbayListingsBucketName, warehouseItem.EbayListingIDs[0], &ebayItemInternal)
	if err != nil {
		return err
	}

	label, err := labels.Create102x152mmItemLabel(
		warehouseItem.ControlAltRepeatID,
		ebayItemInternal.Title,
		strings.Join([]string{"https://www.ebay.co.uk/itm", ebayItemInternal.ID}, "/"),
	)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("102x152-%s.png", warehouseItem.ControlAltRepeatID)

	return labels.UploadFileFromBytes(label, key)
}
