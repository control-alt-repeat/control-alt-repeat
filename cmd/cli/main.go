package main

import (
	"fmt"
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

	cmdEbayImportListing.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	if err := cmdEbayImportListing.MarkFlagRequired("ebay-listing-id"); err != nil {
		log.Fatal().Err(err).Msg("")
		return
	}

	cmdEbayGetNotificationUsage.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	if err := cmdEbayImportListing.MarkFlagRequired("ebay-listing-id"); err != nil {
		log.Fatal().Err(err).Msg("")
		return
	}

	cmdItemMove.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemMove.Flags().StringVar(&shelf, "shelf", "", "Shelf location")
	if err := cmdItemMove.MarkFlagRequired("item-id"); err != nil {
		log.Fatal().Err(err).Msg("")
		return
	}
	if err := cmdItemMove.MarkFlagRequired("shelf"); err != nil {
		log.Fatal().Err(err).Msg("")
		return
	}

	cmdItemRefresh.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	cmdItemRefresh.Flags().BoolVar(&all, "all", false, "All items")
	cmdItemRefresh.MarkFlagsOneRequired("item-id", "all")

	cmdItemPrintShelfLabel.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	if err := cmdItemPrintShelfLabel.MarkFlagRequired("item-id"); err != nil {
		log.Fatal().Err(err).Msg("")
		return
	}

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
