package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal/reconcile"
	"github.com/control-alt-repeat/control-alt-repeat/internal/stripe"
	"github.com/control-alt-repeat/control-alt-repeat/internal/xero"
)

var (
	reconcileStatement        string
	reconcileStatementProfile string
	reconcileProfiles         string
	reconcileRules            string
	reconcilePurchases        []string
	reconcileFrom             string
	reconcileTo               string
	reconcileBankAccount      string
	reconcileLedger           string
	reconcilePlan             string
	reconcileCommit           bool
	reconcileStripe           bool
)

// Reconcile command: "car reconcile"
var cmdReconcile = &cobra.Command{
	Use:   "reconcile",
	Short: "Code bank statement lines for Xero from rules and marketplace order exports",
}

// Plan subcommand: "car reconcile plan"
var cmdReconcilePlan = &cobra.Command{
	Use:   "plan",
	Short: "Builds a reviewable plan CSV from a bank statement - nothing is sent to Xero",
	Run:   reconcilePlanRun,
}

// Apply subcommand: "car reconcile apply"
var cmdReconcileApply = &cobra.Command{
	Use:   "apply",
	Short: "Posts the ready rows of a plan to Xero (dry run unless --commit)",
	Run:   reconcileApplyRun,
}

func registerReconcileCommands() {
	defaultLedger := "xero-ledger.jsonl"
	if dir, err := os.UserConfigDir(); err == nil {
		defaultLedger = filepath.Join(dir, "car", "xero-ledger.jsonl")
	}

	f := cmdReconcilePlan.Flags()
	f.StringVar(&reconcileStatement, "statement", "", "Bank statement CSV")
	f.StringVar(&reconcileStatementProfile, "statement-profile", "xero", "Statement profile name from the profiles file")
	f.StringVar(&reconcileProfiles, "profiles", "configs/reconcile/profiles.json", "CSV profiles file")
	f.StringVar(&reconcileRules, "rules", "configs/reconcile/rules.json", "Coding rules file")
	f.StringArrayVar(&reconcilePurchases, "purchases", nil, "Marketplace order export as profile=path, repeatable, e.g. amazon=orders.csv")
	f.StringVar(&reconcileFrom, "from", "", "Only include statement lines on or after this date (yyyy-mm-dd)")
	f.StringVar(&reconcileTo, "to", "", "Only include statement lines on or before this date (yyyy-mm-dd)")
	f.StringVar(&reconcileBankAccount, "bank-account", "", "Xero bank account (code, name or id) - when set, lines already in Xero are marked 'exists' and unpaid bills are matched")
	f.StringVar(&reconcileLedger, "ledger", defaultLedger, "Ledger of lines already posted by this tool")
	f.StringVar(&reconcilePlan, "out", "plan.csv", "Where to write the plan")
	f.BoolVar(&reconcileStripe, "stripe", false, "Split Stripe payouts into sales, refunds and fees using the Stripe API (STRIPE_API_KEY)")
	_ = cmdReconcilePlan.MarkFlagRequired("statement")

	f = cmdReconcileApply.Flags()
	f.StringVar(&reconcilePlan, "plan", "plan.csv", "Plan CSV produced by `car reconcile plan` and reviewed")
	f.StringVar(&reconcileBankAccount, "bank-account", "", "Xero bank account (code, name or id)")
	f.StringVar(&reconcileLedger, "ledger", defaultLedger, "Ledger of lines already posted by this tool")
	f.BoolVar(&reconcileCommit, "commit", false, "Actually create the transactions in Xero")
	_ = cmdReconcileApply.MarkFlagRequired("bank-account")

	cmdReconcile.AddCommand(cmdReconcilePlan)
	cmdReconcile.AddCommand(cmdReconcileApply)
	cmdRoot.AddCommand(cmdReconcile)
}

func reconcilePlanRun(cmd *cobra.Command, args []string) {
	profiles, err := reconcile.LoadProfiles(reconcileProfiles)
	if err != nil {
		handleError(err)
	}
	rules, err := reconcile.LoadRules(reconcileRules)
	if err != nil {
		handleError(err)
	}

	sp, err := profiles.Statement(reconcileStatementProfile)
	if err != nil {
		handleError(err)
	}
	lines, err := reconcile.LoadStatement(reconcileStatement, sp)
	if err != nil {
		handleError(err)
	}
	lines, err = filterLines(lines, reconcileFrom, reconcileTo)
	if err != nil {
		handleError(err)
	}
	if len(lines) == 0 {
		handleError(fmt.Errorf("no statement lines in range"))
	}
	fmt.Printf("Loaded %d statement lines from %s\n", len(lines), reconcileStatement)

	var sources []reconcile.PurchaseSource
	for _, spec := range reconcilePurchases {
		name, path, ok := strings.Cut(spec, "=")
		if !ok {
			handleError(fmt.Errorf("--purchases must be profile=path, got %q", spec))
		}
		pp, err := profiles.Purchase(name)
		if err != nil {
			handleError(err)
		}
		purchases, err := reconcile.LoadPurchases(path, pp)
		if err != nil {
			handleError(err)
		}
		fmt.Printf("Loaded %d %s orders from %s\n", len(purchases), pp.Source, path)
		sources = append(sources, reconcile.PurchaseSource{Profile: pp, Purchases: purchases})
	}

	ledger, err := reconcile.OpenLedger(reconcileLedger)
	if err != nil {
		handleError(err)
	}

	in := reconcile.PlanInput{Lines: lines, Rules: rules, Sources: sources, Applied: ledger.Applied}

	if reconcileBankAccount != "" {
		c, err := xero.NewClientFromEnv(cmd.Context())
		if err != nil {
			handleError(err)
		}
		account, err := c.FindBankAccount(cmd.Context(), reconcileBankAccount)
		if err != nil {
			handleError(err)
		}
		from, to := lines[0].Date, lines[0].Date
		for _, l := range lines {
			if l.Date.Before(from) {
				from = l.Date
			}
			if l.Date.After(to) {
				to = l.Date
			}
		}
		existing, err := c.ListAccountActivity(cmd.Context(), account.AccountID, from.AddDate(0, 0, -3), to.AddDate(0, 0, 3))
		if err != nil {
			handleError(err)
		}
		fmt.Printf("Found %d existing Xero transactions and payments on %s\n", len(existing), account.Name)
		in.Existing = existing

		bills, err := c.ListUnpaidBills(cmd.Context())
		if err != nil {
			handleError(err)
		}
		fmt.Printf("Found %d unpaid bills\n", len(bills))
		in.Bills = bills
	}

	if reconcileStripe {
		payouts, err := loadStripePayouts(cmd.Context(), lines)
		if err != nil {
			handleError(err)
		}
		in.Payouts = payouts
	}

	rows, err := reconcile.BuildPlan(in)
	if err != nil {
		handleError(err)
	}
	if err := reconcile.WritePlan(reconcilePlan, rows); err != nil {
		handleError(err)
	}

	printSummary(rows)
	fmt.Printf("\nWrote %s - review it, set status to ready/skip, then run `car reconcile apply`\n", reconcilePlan)
}

func reconcileApplyRun(cmd *cobra.Command, args []string) {
	rows, err := reconcile.ReadPlan(reconcilePlan)
	if err != nil {
		handleError(err)
	}
	ledger, err := reconcile.OpenLedger(reconcileLedger)
	if err != nil {
		handleError(err)
	}
	c, err := xero.NewClientFromEnv(cmd.Context())
	if err != nil {
		handleError(err)
	}
	account, err := c.FindBankAccount(cmd.Context(), reconcileBankAccount)
	if err != nil {
		handleError(err)
	}
	fmt.Printf("Bank account: %s (%s)\n", account.Name, account.AccountID)

	res, applyErr := reconcile.Apply(cmd.Context(), c, ledger, rows, reconcile.ApplyOptions{
		BankAccountID: account.AccountID,
		Commit:        reconcileCommit,
		Log:           func(format string, a ...any) { fmt.Printf(format+"\n", a...) },
	})

	// Always save progress, even after an error part way through.
	if reconcileCommit || res.Duplicates > 0 {
		if err := reconcile.WritePlan(reconcilePlan, rows); err != nil {
			handleError(err)
		}
	}
	if applyErr != nil {
		handleError(applyErr)
	}

	fmt.Printf("\nPosted %d, failed %d, already in Xero %d", res.Posted, res.Failed, res.Duplicates)
	if !reconcileCommit {
		fmt.Printf(", would post %d (dry run - add --commit to send)", res.DryRun)
	}
	fmt.Println()
	if res.Posted > 0 {
		fmt.Println("In Xero, open the bank account's Reconcile tab and click OK on each suggested match.")
	}
}

func loadStripePayouts(ctx context.Context, lines []reconcile.StatementLine) ([]reconcile.Payout, error) {
	c, err := stripe.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	from, to := lines[0].Date, lines[0].Date
	for _, l := range lines {
		if l.Date.Before(from) {
			from = l.Date
		}
		if l.Date.After(to) {
			to = l.Date
		}
	}
	payouts, err := c.ListPayouts(ctx, from.AddDate(0, 0, -7), to.AddDate(0, 0, 7))
	if err != nil {
		return nil, err
	}
	out := make([]reconcile.Payout, 0, len(payouts))
	for _, p := range payouts {
		txns, err := c.PayoutTransactions(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, reconcile.StripePayout(p, txns))
	}
	fmt.Printf("Loaded %d Stripe payouts\n", len(out))
	return out, nil
}

func filterLines(lines []reconcile.StatementLine, from, to string) ([]reconcile.StatementLine, error) {
	var start, end time.Time
	var err error
	if from != "" {
		if start, err = time.Parse("2006-01-02", from); err != nil {
			return nil, err
		}
	}
	if to != "" {
		if end, err = time.Parse("2006-01-02", to); err != nil {
			return nil, err
		}
	}
	out := lines[:0:0]
	for _, l := range lines {
		if (from == "" || !l.Date.Before(start)) && (to == "" || !l.Date.After(end)) {
			out = append(out, l)
		}
	}
	return out, nil
}

func printSummary(rows []reconcile.PlanRow) {
	counts := map[string]int{}
	seen := map[string]bool{}
	for _, g := range reconcile.GroupPlan(rows) {
		if seen[g.LineID] {
			continue
		}
		seen[g.LineID] = true
		counts[g.Rows[0].Status]++
	}
	fmt.Println("\nStatement lines by status:")
	for _, s := range []string{reconcile.StatusReady, reconcile.StatusReview, reconcile.StatusExists, reconcile.StatusApplied, reconcile.StatusSkip} {
		if counts[s] > 0 {
			fmt.Printf("  %-8s %d\n", s, counts[s])
		}
	}
}
