package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/spf13/cobra"
)

// Item command: "car item"
var cmdItem = &cobra.Command{
	Use:   "item",
	Short: "Item related operations",
}

// Move item subcommand: "car item move"
var cmdItemMove = &cobra.Command{
	Use:   "move",
	Short: "Move an item to a new shelf",
	Run:   itemMove,
}

// Import listing subcommand: "car item refresh"
var cmdItemRefresh = &cobra.Command{
	Use:   "refresh",
	Short: "Refreshes internal item from source",
	Run:   itemRefresh,
}

// Import listing subcommand: "car item print-shelf-label"
var cmdItemPrintShelfLabel = &cobra.Command{
	Use:   "print-shelf-label",
	Short: "Prints a shelf label for the item",
	Run:   itemPrintShelfLabel,
}

func itemMove(cmd *cobra.Command, args []string) {
	matched, err := regexp.MatchString(`^[A-Z]{3}-[0-9]{3}$`, itemID)
	if err != nil || !matched {
		fmt.Println("Error: item ID must be in the format A-Z-0-9 (e.g., A-123)")
		os.Exit(1)
	}
	fmt.Println("Moving item with ID:", itemID, "to shelf:", shelf)

	err = internal.MoveItem(itemID, shelf)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func itemRefresh(cmd *cobra.Command, args []string) {
	// Parse flags
	all, _ := cmd.Flags().GetBool("all")
	itemID, _ := cmd.Flags().GetString("item-id")

	// Check for conflicting or missing flags
	if all && itemID != "" {
		fmt.Println("Error: You cannot specify both --all and --item-id.")
		os.Exit(1)
	}

	if !all && itemID == "" {
		fmt.Println("Error: You must specify either --all or --item-id.")
		os.Exit(1)
	}

	var err error
	if all {
		err = internal.RefreshItemsFromEbay()
	} else if itemID != "" {
		err = internal.RefreshItemFromEbay(itemID)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func itemPrintShelfLabel(cmd *cobra.Command, args []string) {
	matched, err := regexp.MatchString(`^[A-Z]{3}-[0-9]{3}$`, itemID)
	if err != nil || !matched {
		fmt.Println("Error: item ID must be in the format A-Z-0-9 (e.g., A-123)")
		os.Exit(1)
	}
	fmt.Println("Printing shelf label for item ID:", itemID)

	err = internal.ItemPrintShelfLabel(itemID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
