package reconcile

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/shopspring/decimal"
)

// Rules is the reusable, human-maintained knowledge of how to code each kind of
// bank line. Rules are evaluated top to bottom; the first match wins.
type Rules struct {
	// Used for purchase items that no item rule matches.
	DefaultPurchaseAccountCode string `json:"default_purchase_account_code"`
	DefaultPurchaseTaxType     string `json:"default_purchase_tax_type,omitempty"`

	Rules []Rule `json:"rules"`
}

type Rule struct {
	Name string `json:"name"`

	// Conditions - all that are set must match. Patterns are Go regexps, use (?i) for case-insensitive.
	Text      string           `json:"text,omitempty"`      // payee + description + reference
	Item      string           `json:"item,omitempty"`      // purchase item description (item rules only)
	Direction string           `json:"direction,omitempty"` // "in", "out" or empty for either
	MinAmount *decimal.Decimal `json:"min_amount,omitempty"`
	MaxAmount *decimal.Decimal `json:"max_amount,omitempty"`

	// Coding applied in Xero.
	Contact     string `json:"contact,omitempty"`
	AccountCode string `json:"account_code"`
	TaxType     string `json:"tax_type,omitempty"`
	Narrative   string `json:"narrative,omitempty"`
	// Leave the line for a human even though the coding is known (e.g. mixed-use cards).
	Review bool `json:"review,omitempty"`
	// Never post this line, e.g. marketplace payouts Link My Books already accounts for.
	Skip bool `json:"skip,omitempty"`

	textRe *regexp.Regexp
	itemRe *regexp.Regexp
}

func LoadRules(path string) (*Rules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r Rules
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	for i := range r.Rules {
		rule := &r.Rules[i]
		if rule.Name == "" {
			return nil, fmt.Errorf("%s: rule %d has no name", path, i+1)
		}
		if rule.AccountCode == "" && !rule.Skip {
			return nil, fmt.Errorf("%s: rule %q has no account_code", path, rule.Name)
		}
		if rule.Direction != "" && rule.Direction != "in" && rule.Direction != "out" {
			return nil, fmt.Errorf("%s: rule %q direction must be in, out or empty", path, rule.Name)
		}
		if rule.Text != "" {
			if rule.textRe, err = regexp.Compile(rule.Text); err != nil {
				return nil, fmt.Errorf("%s: rule %q text: %w", path, rule.Name, err)
			}
		}
		if rule.Item != "" {
			if rule.itemRe, err = regexp.Compile(rule.Item); err != nil {
				return nil, fmt.Errorf("%s: rule %q item: %w", path, rule.Name, err)
			}
		}
	}
	return &r, nil
}

// MatchLine finds the rule for a whole bank line (item rules are ignored).
func (r *Rules) MatchLine(l StatementLine) *Rule {
	for i := range r.Rules {
		rule := &r.Rules[i]
		if rule.itemRe != nil {
			continue
		}
		if rule.matchesLine(l) {
			return rule
		}
	}
	return nil
}

// MatchItem finds the rule for one purchase item within a bank line.
func (r *Rules) MatchItem(l StatementLine, item PurchaseItem) *Rule {
	for i := range r.Rules {
		rule := &r.Rules[i]
		if rule.itemRe == nil || !rule.itemRe.MatchString(item.Description) {
			continue
		}
		if rule.matchesLine(l) {
			return rule
		}
	}
	return nil
}

func (rule *Rule) matchesLine(l StatementLine) bool {
	if rule.Direction == "in" && l.IsSpend() || rule.Direction == "out" && !l.IsSpend() {
		return false
	}
	abs := l.Amount.Abs()
	if rule.MinAmount != nil && abs.LessThan(*rule.MinAmount) {
		return false
	}
	if rule.MaxAmount != nil && abs.GreaterThan(*rule.MaxAmount) {
		return false
	}
	if rule.textRe != nil && !rule.textRe.MatchString(l.Text()) {
		return false
	}
	return true
}
