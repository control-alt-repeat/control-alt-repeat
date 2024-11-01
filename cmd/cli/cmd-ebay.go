package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
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

	_, err := internal.ImportEbayListingByID(cmd.Context(), ebayListingID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebayGetNotificationUsage(cmd *cobra.Command, args []string) {
	err := ebay.GetNotificationUsage(cmd.Context(), ebayListingID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ebaySetNotificationPreferences(cmd *cobra.Command, args []string) {
	err := ebay.SetNotificationPreferences(cmd.Context())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
