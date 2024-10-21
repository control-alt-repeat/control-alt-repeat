package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/spf13/cobra"
)

// eBay command: "car ebay"
var cmdEbay = &cobra.Command{
	Use:   "ebay",
	Short: "eBay related operations",
}

// Import listing subcommand: "car ebay import-listing"
var cmdEbayImportListing = &cobra.Command{
	Use:   "import-listing",
	Short: "Import a listing from eBay",
	Run:   ebayImportListing,
}

// Import listing subcommand: "car ebay get-notification-preferences"
var cmdEbayGetNotificationUsage = &cobra.Command{
	Use:   "get-notification-usage",
	Short: "Get notification usage from eBay",
	Run:   ebayGetNotificationUsage,
}

// Import listing subcommand: "car ebay set-notification-preferences"
var cmdEbaySetNotificationPreferences = &cobra.Command{
	Use:   "set-notification-preferences",
	Short: "Set notification preferences on eBay",
	Run:   ebaySetNotificationPreferences,
}

func ebayImportListing(cmd *cobra.Command, args []string) {
	if _, err := strconv.Atoi(ebayListingID); err != nil {
		fmt.Println("Error: eBay listing ID must be a valid integer")
		os.Exit(1)
	}
	fmt.Println("Importing eBay listing with ID:", ebayListingID)

	err := internal.ImportEbayListingByID(ebayListingID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebayGetNotificationUsage(cmd *cobra.Command, args []string) {
	err := ebay.GetNotificationUsage(ebayListingID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebaySetNotificationPreferences(cmd *cobra.Command, args []string) {
	err := ebay.SetNotificationPreferences()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
