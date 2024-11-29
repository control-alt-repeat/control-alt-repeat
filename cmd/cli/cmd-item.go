package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
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

	err = internal.MoveItem(cmd.Context(), itemID, shelf)

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

	err = internal.ItemPrintShelfLabel(cmd.Context(), internal.ItemPrintShelfLabelOptions{
		ItemID: itemID,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
