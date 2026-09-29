# Xero reconciliation

Repeatable, rule-driven coding of bank statement lines for Xero. No AI in the loop:
the same inputs always give the same plan.

```
bank CSV ─┐
orders ───┼─> car reconcile plan ─> plan.csv ─> (you review) ─> car reconcile apply --commit ─> Xero
rules ────┘                                                                                      │
                                                         Xero Reconcile tab: click OK on each match
```

## One-off setup

1. Create a Xero app at https://developer.xero.com/app/manage
   - **Custom connection** (simplest, one organisation, small monthly fee): no login step.
   - **Web app** (free): redirect URI `http://localhost:8765/callback`, then run `car xero login` once.
2. Export credentials (never commit them):
   ```sh
   export XERO_CLIENT_ID=...
   export XERO_CLIENT_SECRET=...
   ```
   Tokens are cached in `~/.config/car/xero-token.json` (mode 600).
3. `car xero accounts` - lists account codes and tax types.
4. `cp configs/reconcile/rules.example.json configs/reconcile/rules.json` and edit it using those codes.

## Each time

Put inputs in `reconcile-data/` (git-ignored).

```sh
go run ./cmd/cli reconcile plan \
  --statement reconcile-data/bank-2026-03.csv \
  --purchases amazon=reconcile-data/amazon-orders.csv \
  --bank-account "Business Current Account" \
  --out reconcile-data/plan-2026-03.csv

# open the plan, fix/complete rows marked review, set them to ready or skip

go run ./cmd/cli reconcile apply --plan reconcile-data/plan-2026-03.csv --bank-account "Business Current Account"          # dry run
go run ./cmd/cli reconcile apply --plan reconcile-data/plan-2026-03.csv --bank-account "Business Current Account" --commit
```

Then in Xero open the bank account's **Reconcile** tab: each posted transaction appears
as the suggested match for its statement line - click **OK**.

When you find yourself coding the same thing by hand twice, add a rule instead.

## Plan statuses

| status  | meaning |
|---------|---------|
| ready   | fully coded; `apply` posts it |
| review  | needs a decision - `reason` says why. Fill in contact/account_code/tax_type and set to `ready`, or `skip` |
| skip    | not handled by this tool (e.g. transfers between your own accounts - do those in Xero) |
| exists  | a transaction with the same amount (±3 days) is already in Xero |
| applied | posted by this tool; `xero_id` is the Xero BankTransactionID |

A bank line split over several order items has one row per item sharing a `line_id`;
their `line_amount`s must add up to the bank amount or `apply` refuses to run.

## Safety

- `plan` never writes to Xero.
- `apply` is a dry run unless `--commit`.
- Transactions are created **unreconciled** so the bank statement stays the source of truth.
- Every posted line is recorded in `~/.config/car/xero-ledger.jsonl`; a line is never posted twice,
  and Xero is re-checked for a matching transaction immediately before posting.
- A charge that could match more than one order is sent to review, never guessed.

## Files

- `profiles.json` - how to read each bank / marketplace CSV. New bank or shop = new profile, no code.
- `rules.example.json` - coding rules; first match wins. `text` matches payee+description+reference,
  `item` matches purchase item descriptions, `review: true` forces a human look.
