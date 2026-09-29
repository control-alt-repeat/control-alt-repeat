package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal/xero"
)

// Xero command: "car xero"
var cmdXero = &cobra.Command{
	Use:   "xero",
	Short: "Xero related operations (credentials come from XERO_* environment variables)",
}

// Login subcommand: "car xero login" - only needed for Xero "Web app" credentials
var cmdXeroLogin = &cobra.Command{
	Use:   "login",
	Short: "Authorise a Xero Web app (not needed for Custom connections)",
	Run:   xeroLogin,
}

// Connections subcommand: "car xero connections"
var cmdXeroConnections = &cobra.Command{
	Use:   "connections",
	Short: "Lists the Xero organisations these credentials can access",
	Run:   xeroConnections,
}

// Accounts subcommand: "car xero accounts"
var cmdXeroAccounts = &cobra.Command{
	Use:   "accounts",
	Short: "Lists the chart of accounts and tax rates, for writing reconcile rules",
	Run:   xeroAccounts,
}

func registerXeroCommands() {
	cmdXero.AddCommand(cmdXeroLogin)
	cmdXero.AddCommand(cmdXeroConnections)
	cmdXero.AddCommand(cmdXeroAccounts)
	cmdRoot.AddCommand(cmdXero)
}

func xeroLogin(cmd *cobra.Command, args []string) {
	cfg, err := xero.ConfigFromEnv()
	if err != nil {
		handleError(err)
	}
	if err := xero.Login(cmd.Context(), cfg, os.Stdin, os.Stdout); err != nil {
		handleError(err)
	}
	fmt.Println("Saved token to", cfg.TokenFile)
}

func xeroConnections(cmd *cobra.Command, args []string) {
	cfg, err := xero.ConfigFromEnv()
	if err != nil {
		handleError(err)
	}
	ts, err := xero.NewTokenSource(cfg)
	if err != nil {
		handleError(err)
	}
	conns, err := xero.ListConnections(cmd.Context(), ts)
	if err != nil {
		handleError(err)
	}
	for _, c := range conns {
		fmt.Printf("%s  %s  %s\n", c.TenantID, c.TenantType, c.TenantName)
	}
}

func xeroAccounts(cmd *cobra.Command, args []string) {
	c, err := xero.NewClientFromEnv(cmd.Context())
	if err != nil {
		handleError(err)
	}

	accounts, err := c.ListAccounts(cmd.Context())
	if err != nil {
		handleError(err)
	}
	sort.Slice(accounts, func(i, j int) bool {
		if accounts[i].Type != accounts[j].Type {
			return accounts[i].Type < accounts[j].Type
		}
		return accounts[i].Code < accounts[j].Code
	})

	fmt.Println("ACCOUNTS (code, type, default tax, name, id)")
	for _, a := range accounts {
		if a.Status != "ACTIVE" {
			continue
		}
		fmt.Printf("  %-6s %-12s %-14s %-40s %s\n", a.Code, a.Type, a.TaxType, a.Name, a.AccountID)
	}

	rates, err := c.ListTaxRates(cmd.Context())
	if err != nil {
		handleError(err)
	}
	fmt.Println("\nTAX RATES (tax_type, rate, name)")
	for _, r := range rates {
		if r.Status != "ACTIVE" {
			continue
		}
		fmt.Printf("  %-16s %6s%%  %s\n", r.TaxType, r.EffectiveRate.String(), r.Name)
	}
}
