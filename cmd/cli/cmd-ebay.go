package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/traditionalapi"
)

var cmdEbay = &cobra.Command{
	Use:   "ebay",
	Short: "eBay related operations",
}

var cmdEbayImportListing = &cobra.Command{
	Use:   "import-listing",
	Short: "Import a listing from eBay",
	Run:   ebayImportListing,
}

var cmdEbayGetNotificationUsage = &cobra.Command{
	Use:   "get-notification-usage",
	Short: "Get notification usage from eBay",
	Run:   ebayGetNotificationUsage,
}

var cmdEbaySetNotificationPreferences = &cobra.Command{
	Use:   "set-notification-preferences",
	Short: "Set notification preferences on eBay",
	Run:   ebaySetNotificationPreferences,
}

var cmdEbayInventorySetup = &cobra.Command{
	Use:   "inventory-setup",
	Short: "Setting up inventory on eBay",
	Run:   ebayInventorySetup,
}

var cmdEbayInventoryImportListing = &cobra.Command{
	Use:   "inventory-import-listing",
	Short: "Import a listing into inventory on eBay",
	Run:   ebayInventoryImportListing,
}

func ebayImportListing(cmd *cobra.Command, args []string) {
	if _, err := strconv.Atoi(ebayListingID); err != nil {
		fmt.Println("Error: eBay listing ID must be a valid integer")
		os.Exit(1)
	}
	fmt.Println("Importing eBay listing with ID:", ebayListingID)

	_, err := internal.ImportEbayListingByID(cmd.Context(), ebayListingID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebayGetNotificationUsage(cmd *cobra.Command, args []string) {
	err := traditionalapi.GetNotificationUsage(cmd.Context(), ebayListingID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebaySetNotificationPreferences(cmd *cobra.Command, args []string) {
	err := traditionalapi.SetNotificationPreferences(cmd.Context())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebayInventorySetup(cmd *cobra.Command, args []string) {
	err := ebay.SetControlAltRepeatWorkshopLocation(cmd.Context())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebayInventoryImportListing(cmd *cobra.Command, args []string) {
	err := ebay.InventoryImportListing(cmd.Context(), ebayListingID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
