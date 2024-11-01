package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/spf13/cobra"
)

var log zerolog.Logger

var (
	ebayListingID string
	itemID        string
	shelf         string
	contactID     string
	all           bool
)

// Root command: "car"
var cmdRoot = &cobra.Command{
	Use:   "car",
	Short: "Car CLI for managing eBay listings and items",
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	log = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Str("service", "cli").
		Logger().
		Level(zerolog.DebugLevel)

	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	log = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Str("service", "cli").
		Logger().
		Level(zerolog.DebugLevel)

	cmdRoot.AddCommand(cmdEbay)

	cmdEbayImportListing.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	if err := cmdEbayImportListing.MarkFlagRequired("ebay-listing-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdEbay.AddCommand(cmdEbayImportListing)

	cmdEbayGetNotificationUsage.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	if err := cmdEbayGetNotificationUsage.MarkFlagRequired("ebay-listing-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdEbay.AddCommand(cmdEbayGetNotificationUsage)
	cmdEbay.AddCommand(cmdEbaySetNotificationPreferences)

	cmdRoot.AddCommand(cmdFreeagent)

	cmdFreeagentGetContact.Flags().StringVar(&contactID, "contact-id", "", "Freeagent contact ID")
	if err := cmdFreeagentGetContact.MarkFlagRequired("contact-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdFreeagent.AddCommand(cmdFreeagentGetContact)
	cmdFreeagent.AddCommand(cmdFreeagentListContacts)

	cmdRoot.AddCommand(cmdItem)

	cmdItemMove.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemMove.Flags().StringVar(&shelf, "shelf", "", "Shelf location")
	if err := cmdItemMove.MarkFlagRequired("item-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	if err := cmdItemMove.MarkFlagRequired("shelf"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdItem.AddCommand(cmdItemMove)

	cmdItemRefresh.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemRefresh.Flags().BoolVar(&all, "all", false, "All items")
	cmdItemRefresh.MarkFlagsOneRequired("item-id", "all")
	cmdItem.AddCommand(cmdItemRefresh)

	cmdItemPrintShelfLabel.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	if err := cmdItemPrintShelfLabel.MarkFlagRequired("item-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdItem.AddCommand(cmdItemPrintShelfLabel)
	cmdItem.AddCommand(cmdItemMigrateToDynamoDB)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cmdRoot.ExecuteContext(ctx); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}

	return 0
}
