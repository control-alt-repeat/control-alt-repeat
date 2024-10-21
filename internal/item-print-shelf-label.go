package internal

import (
	"fmt"
	"strings"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/labels"
)

func ItemPrintShelfLabel(itemID string) error {
	var warehouseItem WarehouseItem
	var ebayItemInternal EbayItemInternal

	fmt.Println("loading warehouse item ", itemID)
	err := aws.LoadJsonObjectS3(WarehouseItemsBucketName, itemID, &warehouseItem)
	if err != nil {
		return err
	}

	fmt.Println("loading ebay item ", warehouseItem.EbayListingIDs[0])
	err = aws.LoadJsonObjectS3(EbayListingsBucketName, warehouseItem.EbayListingIDs[0], &ebayItemInternal)
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

	fmt.Println("Saving label to S3 buffer: ", key)
	err = aws.SaveBytesToS3("control-alt-repeat-label-print-buffer", key, label, "image/png")
	if err != nil {
		return err
	}

	return labels.NotifyLabelPrintServer()
}
