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
