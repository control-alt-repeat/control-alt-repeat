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
	file          string
	shelf         string
	contactID     string
	closingFunds  string
	fromDate      string
	toDate        string
	bankAccountID string
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

	cmdRoot.AddCommand(cmdAccounting)

	cmdAccountingExplainEbay.Flags().StringVar(&file, "file", "", "eBay transactions CSV")
	if err := cmdAccountingExplainEbay.MarkFlagRequired("file"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdAccountingExplainEbay.Flags().StringVar(&closingFunds, "closing-funds", "", "Expected closing funds, used to verify before upload, e.g. 30.55")
	if err := cmdAccountingExplainEbay.MarkFlagRequired("closing-funds"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdAccountingExplainEbay.Flags().StringVar(&fromDate, "from-date", "", "Date to start explaining")
	if err := cmdAccountingExplainEbay.MarkFlagRequired("from-date"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdAccountingExplainEbay.Flags().StringVar(&toDate, "to-date", "", "Date to start explaining")
	if err := cmdAccountingExplainEbay.MarkFlagRequired("to-date"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdAccountingExplainEbay.Flags().StringVar(&bankAccountID, "bank-account-id", "", "FreeAgent bank account ID to explain transactions")
	if err := cmdAccountingExplainEbay.MarkFlagRequired("bank-account-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdAccounting.AddCommand(cmdAccountingExplainEbay)

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
	cmdEbay.AddCommand(cmdEbayInventorySetup)

	cmdEbayInventoryImportListing.Flags().StringVar(&ebayListingID, "ebay-listing-id", "", "eBay listing ID")
	if err := cmdEbayInventoryImportListing.MarkFlagRequired("ebay-listing-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdEbay.AddCommand(cmdEbayInventoryImportListing)

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

	cmdItemPrintShelfLabel.Flags().StringVar(&itemID, "item-id", "", "Item ID")
	if err := cmdItemPrintShelfLabel.MarkFlagRequired("item-id"); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}
	cmdItem.AddCommand(cmdItemPrintShelfLabel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cmdRoot.ExecuteContext(ctx); err != nil {
		log.Error().Err(err).Msg("")
		return 1
	}

	return 0
}
