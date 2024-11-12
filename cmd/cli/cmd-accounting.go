package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/reports"
)

// Accounting command: "car accounting"
var cmdAccounting = &cobra.Command{
	Use:   "accounting",
	Short: "Accounting related operations",
}

// Move item subcommand: "car item move"
var cmdAccountingExplainEbay = &cobra.Command{
	Use:   "explain-ebay",
	Short: "Add eBay transactions to accounting software",
	Run:   accountingExplainEbay,
}

func accountingExplainEbay(cmd *cobra.Command, args []string) {
	fmt.Println("Loading eBay transcation CSV:", file)

	_, err := reports.LoadTransactionsFile(cmd.Context(), file)

	// fmt.Println(report.Transactions[7].Description)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
