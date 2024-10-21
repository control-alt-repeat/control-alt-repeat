package internal

import (
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

func RefreshItemsFromEbay() error {
	return aws.IterateS3Objects(WarehouseItemsBucketName, "eu-west-2", RefreshItemFromEbay)
}

func RefreshItemFromEbay(itemID string) error {
	var warehouseItem WarehouseItem
	var ebayItemInternal EbayItemInternal

	fmt.Println("loading warehouse item ", itemID)
	err := aws.LoadJsonObjectS3(WarehouseItemsBucketName, itemID, &warehouseItem)
	if err != nil {
		return err
	}

	fmt.Println("loading ebay item ", warehouseItem.Ebay.ID)
	err = aws.LoadJsonObjectS3(EbayListingsBucketName, warehouseItem.Ebay.ID, &ebayItemInternal)
	if err != nil {
		return err
	}

	ebayItemSource, err := ebay.GetItem(warehouseItem.Ebay.ID, []string{
		"ItemID",
		"Title",
		"PictureDetails",
		"ListingDetails",
	})
	if err != nil {
		return err
	}

	ebayItemInternal.Title = ebayItemSource.Title
	ebayItemInternal.PictureURL = ebayItemSource.PictureDetails.PictureURL[0]
	ebayItemInternal.ViewItemURL = ebayItemSource.ListingDetails.ViewItemURL

	fmt.Println("Title: ", ebayItemInternal.Title)
	fmt.Println("PictureURL: ", ebayItemInternal.PictureURL)
	fmt.Println("ViewItemURL: ", ebayItemInternal.ViewItemURL)

	return aws.SaveJsonObjectS3(EbayListingsBucketName, warehouseItem.Ebay.ID, &ebayItemInternal)
}
