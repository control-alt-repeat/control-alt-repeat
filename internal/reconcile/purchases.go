package reconcile

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Purchase is one order from a marketplace (Amazon, eBay, ...), used to itemise
// the matching bank line in Xero.
type Purchase struct {
	Source  string
	OrderID string
	Date    time.Time
	Seller  string
	URL     string
	Items   []PurchaseItem
}

type PurchaseItem struct {
	Description string
	Quantity    int
	Total       decimal.Decimal // gross, including VAT and its share of postage
}

func (p Purchase) Total() decimal.Decimal {
	total := decimal.Zero
	for _, i := range p.Items {
		total = total.Add(i.Total)
	}
	return total
}

// PurchaseProfile maps a marketplace order export onto Purchase. Rows sharing an
// order ID are grouped into one Purchase with one item per row.
type PurchaseProfile struct {
	Name            string   `json:"name"`
	Source          string   `json:"source"`
	Contact         string   `json:"contact,omitempty"` // Xero contact, defaults to Source
	SkipRows        int      `json:"skip_rows,omitempty"`
	OrderIDColumn   string   `json:"order_id_column"`
	DateColumn      string   `json:"date_column"`
	DateFormats     []string `json:"date_formats,omitempty"`
	ItemTotalColumn string   `json:"item_total_column"`
	// Optional: what was actually charged for the whole order (repeated on each row).
	// When set it wins over the item sum - item amounts are scaled to fit, so
	// discounts are spread across the items - and £0 orders are dropped.
	OrderTotalColumn  string `json:"order_total_column,omitempty"`
	DescriptionColumn string `json:"description_column"`
	QuantityColumn    string `json:"quantity_column,omitempty"`
	SellerColumn      string `json:"seller_column,omitempty"`
	URLColumn         string `json:"url_column,omitempty"` // link to the order, attached in Xero
	// Rows where ExcludeColumn matches ExcludePattern are dropped (e.g. cancelled orders).
	ExcludeColumn  string `json:"exclude_column,omitempty"`
	ExcludePattern string `json:"exclude_pattern,omitempty"`
	// Bank lines must match this pattern to be considered for this source's orders.
	BankPayeePattern string `json:"bank_payee_pattern"`
	// How many days either side of the order date the card charge may land.
	DaysBefore int `json:"days_before,omitempty"`
	DaysAfter  int `json:"days_after,omitempty"`
}

func LoadPurchases(path string, p PurchaseProfile) ([]Purchase, error) {
	t, err := readTable(path, p.SkipRows)
	if err != nil {
		return nil, err
	}
	if err := t.require(p.OrderIDColumn, p.DateColumn, p.ItemTotalColumn, p.OrderTotalColumn, p.DescriptionColumn, p.QuantityColumn, p.SellerColumn, p.URLColumn, p.ExcludeColumn); err != nil {
		return nil, err
	}

	var exclude *regexp.Regexp
	if p.ExcludePattern != "" {
		if exclude, err = regexp.Compile(p.ExcludePattern); err != nil {
			return nil, fmt.Errorf("profile %q exclude_pattern: %w", p.Name, err)
		}
	}

	byOrder := map[string]*Purchase{}
	charged := map[string]decimal.Decimal{}
	var order []string

	for i, row := range t.rows {
		rowNum := i + 2 + p.SkipRows
		if exclude != nil && exclude.MatchString(t.value(row, p.ExcludeColumn)) {
			continue
		}

		id := t.value(row, p.OrderIDColumn)
		if id == "" {
			return nil, fmt.Errorf("%s row %d: empty order id", path, rowNum)
		}

		total, err := parseMoney(t.value(row, p.ItemTotalColumn))
		if err != nil {
			return nil, fmt.Errorf("%s row %d: %w", path, rowNum, err)
		}
		qty := 1
		if q := t.value(row, p.QuantityColumn); q != "" {
			if qty, err = strconv.Atoi(q); err != nil {
				return nil, fmt.Errorf("%s row %d: quantity %q: %w", path, rowNum, q, err)
			}
		}

		purchase, ok := byOrder[id]
		if !ok {
			date, err := parseDate(t.value(row, p.DateColumn), p.DateFormats)
			if err != nil {
				return nil, fmt.Errorf("%s row %d: %w", path, rowNum, err)
			}
			purchase = &Purchase{Source: p.Source, OrderID: id, Date: date, Seller: t.value(row, p.SellerColumn), URL: t.value(row, p.URLColumn)}
			byOrder[id] = purchase
			order = append(order, id)

			if p.OrderTotalColumn != "" {
				c, err := parseMoney(t.value(row, p.OrderTotalColumn))
				if err != nil {
					return nil, fmt.Errorf("%s row %d: %w", path, rowNum, err)
				}
				charged[id] = c.Abs()
			}
		}
		purchase.Items = append(purchase.Items, PurchaseItem{
			Description: t.value(row, p.DescriptionColumn),
			Quantity:    qty,
			Total:       total.Abs(),
		})
	}

	purchases := make([]Purchase, 0, len(order))
	for _, id := range order {
		purchase := byOrder[id]
		if c, ok := charged[id]; ok {
			if c.IsZero() {
				continue // nothing was charged (cancelled, refunded or paid by voucher)
			}
			purchase.Items = scaleItems(purchase.Items, c)
		}
		purchases = append(purchases, *purchase)
	}
	sort.SliceStable(purchases, func(i, j int) bool { return purchases[i].Date.Before(purchases[j].Date) })
	return purchases, nil
}

// scaleItems makes item amounts add up exactly to the amount charged, spreading any
// difference (discounts, coupons) in proportion. The last item absorbs rounding pennies.
func scaleItems(items []PurchaseItem, charged decimal.Decimal) []PurchaseItem {
	sum := decimal.Zero
	for _, i := range items {
		sum = sum.Add(i.Total)
	}
	if sum.Equal(charged) {
		return items
	}

	out := make([]PurchaseItem, len(items))
	copy(out, items)
	allocated := decimal.Zero
	for i := range out {
		if i == len(out)-1 {
			out[i].Total = charged.Sub(allocated)
			break
		}
		if sum.IsZero() {
			out[i].Total = decimal.Zero
			continue
		}
		out[i].Total = out[i].Total.Mul(charged).Div(sum).Round(2)
		allocated = allocated.Add(out[i].Total)
	}
	return out
}

// contactName is the Xero contact used for this source, e.g. "amazon" -> "Amazon".
func (p PurchaseProfile) contactName() string {
	if p.Contact != "" {
		return p.Contact
	}
	name := p.Source
	if name == "" {
		name = p.Name
	}
	if name == "" {
		return "Unknown supplier"
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func (p PurchaseProfile) payeeMatcher() (*regexp.Regexp, error) {
	pattern := p.BankPayeePattern
	if pattern == "" {
		pattern = "(?i)" + regexp.QuoteMeta(strings.TrimSpace(p.Source))
	}
	return regexp.Compile(pattern)
}
