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

	cmdEbayImportListing.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	if err := cmdEbayImportListing.MarkFlagRequired("ebay-listing-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}

	cmdEbayGetNotificationUsage.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	if err := cmdEbayImportListing.MarkFlagRequired("ebay-listing-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}

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

	cmdItemRefresh.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemRefresh.Flags().BoolVar(&all, "all", false, "All items")
	cmdItemRefresh.MarkFlagsOneRequired("item-id", "all")

	cmdItemPrintShelfLabel.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	if err := cmdItemPrintShelfLabel.MarkFlagRequired("item-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}

	cmdEbay.AddCommand(cmdEbayImportListing)
	cmdEbay.AddCommand(cmdEbayGetNotificationUsage)
	cmdEbay.AddCommand(cmdEbaySetNotificationPreferences)
	cmdItem.AddCommand(cmdItemMove)
	cmdItem.AddCommand(cmdItemRefresh)
	cmdItem.AddCommand(cmdItemPrintShelfLabel)
	cmdRoot.AddCommand(cmdEbay)
	cmdRoot.AddCommand(cmdItem)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cmdRoot.ExecuteContext(ctx); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}

	return 0
}
