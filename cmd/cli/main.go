package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ebayListingID string
	itemID        string
	shelf         string
	all           bool
)

// Root command: "car"
var cmdRoot = &cobra.Command{
	Use:   "car",
	Short: "Car CLI for managing eBay listings and items",
}

func main() {
	cmdEbayImportListing.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	cmdEbayImportListing.MarkFlagRequired("ebay-listing-id")

	cmdEbayGetNotificationUsage.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	cmdEbayImportListing.MarkFlagRequired("ebay-listing-id")

	cmdItemMove.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemMove.Flags().StringVar(&shelf, "shelf", "", "Shelf location")
	cmdItemMove.MarkFlagRequired("item-id")
	cmdItemMove.MarkFlagRequired("shelf")

	cmdItemRefresh.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemRefresh.Flags().BoolVar(&all, "all", false, "All items")
	cmdItemRefresh.MarkFlagsOneRequired("item-id", "all")

	cmdItemPrintShelfLabel.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemPrintShelfLabel.MarkFlagRequired("item-id")

	cmdEbay.AddCommand(cmdEbayImportListing)
	cmdEbay.AddCommand(cmdEbayGetNotificationUsage)
	cmdEbay.AddCommand(cmdEbaySetNotificationPreferences)
	cmdItem.AddCommand(cmdItemMove)
	cmdItem.AddCommand(cmdItemRefresh)
	cmdItem.AddCommand(cmdItemPrintShelfLabel)
	cmdRoot.AddCommand(cmdEbay)
	cmdRoot.AddCommand(cmdItem)

	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
