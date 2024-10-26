package internal

import (
	"context"
	"fmt"
	"strings"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/labels"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/warehouse"
)

func ItemPrintShelfLabel(ctx context.Context, itemID string) error {
	var warehouseItem warehouse.WarehouseItem
	var ebayItemInternal warehouse.EbayItemInternal

	err := aws.LoadJsonObjectS3(ctx, warehouse.WarehouseItemsBucketName, itemID, &warehouseItem)
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
