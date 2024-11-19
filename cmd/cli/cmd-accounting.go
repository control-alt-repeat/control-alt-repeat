package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal/accounting"
	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/reports"
	"github.com/control-alt-repeat/control-alt-repeat/internal/freeagent"
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
	startDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		handleError(err)
	}

	endDate, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		handleError(err)
	}

	expectedClosingFunds, err := decimal.NewFromString(closingFunds)
	if err != nil {
		handleError(err)
	}

	fmt.Println("Loading eBay transcation CSV:", file)

	report, err := reports.LoadTransactionsFile(cmd.Context(), file)
	if err != nil {
		handleError(err)
	}

	var filtered []reports.Transaction
	for _, t := range report.Transactions {
		if !t.TransactionCreationDate.Before(startDate) && !t.TransactionCreationDate.After(endDate) {
			filtered = append(filtered, t)
		}
	}

	explainations := []freeagent.BankTransactionExplanation{}

	for _, transaction := range filtered {
		result, err := accounting.MapEbayTransactionsToFreeAgent(cmd.Context(), transaction)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		explainations = append(explainations, result...)
	}

	sort.Slice(explainations, func(i, j int) bool {
		return explainations[i].DatedOn.ToTime().Before(explainations[j].DatedOn.ToTime())
	})

	bankAccount, err := freeagent.GetBankAccount(cmd.Context(), freeagent.GetBankAccountOptions{
		BankAccountID: bankAccountID,
	})
	if err != nil {
		handleError(err)
	}

	newClosingFunds := bankAccount.CurrentBalance
	for _, e := range explainations {
		grossDecimal, err := decimal.NewFromString(e.GrossValue)
		if err != nil {
			handleError(err)
		}

		newClosingFunds = newClosingFunds.Add(grossDecimal)

		fmt.Printf("%6s %s %6s %s\n", newClosingFunds, e.DatedOn.ToTime().Format("2006-01-02"), e.GrossValue, e.Description)
	}

	if !expectedClosingFunds.Equal(newClosingFunds) {
		handleError(fmt.Errorf("expected closing funds (%s) does not match actual (%s)", expectedClosingFunds, newClosingFunds))
	}

	for _, e := range explainations {
		e.BankAccount = bankAccount.ID

		err = freeagent.CreateBankTransactionExplaination(cmd.Context(), freeagent.CreateBankTransactionExplainationRequest{
			BankTransactionExplanation: e,
		})
		if err != nil {
			handleError(err)
		}

		fmt.Printf("Created explanation: %s %6s %s\n", e.DatedOn.ToTime().Format("2006-01-02"), e.GrossValue, e.Description)
	}

	if err != nil {
		handleError(err)
	}
}
