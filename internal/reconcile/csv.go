package reconcile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var defaultDateFormats = []string{
	"02/01/2006",
	"2/1/2006",
	"2006-01-02",
	"02 Jan 2006",
	"2 Jan 2006",
	"02-Jan-2006",
	"02/01/06",
	time.RFC3339,
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05",
	"Jan 2, 2006",
	"January 2, 2006",
}

// table is a CSV file with case-insensitive column lookup.
type table struct {
	header map[string]int
	rows   [][]string
	path   string
}

func readTable(path string, skipRows int) (*table, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(stripBOM(f))
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	var records [][]string
	for {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		records = append(records, rec)
	}
	if len(records) <= skipRows {
		return nil, fmt.Errorf("%s: no header row found", path)
	}
	records = records[skipRows:]

	t := &table{header: map[string]int{}, path: path}
	for i, h := range records[0] {
		t.header[normaliseHeader(h)] = i
	}
	for _, rec := range records[1:] {
		if isBlank(rec) {
			continue
		}
		t.rows = append(t.rows, rec)
	}
	return t, nil
}

func stripBOM(r io.Reader) io.Reader {
	buf := make([]byte, 3)
	n, _ := io.ReadFull(r, buf)
	if n == 3 && buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		return r
	}
	return io.MultiReader(strings.NewReader(string(buf[:n])), r)
}

func normaliseHeader(h string) string {
	return strings.ToLower(strings.TrimSpace(h))
}

func isBlank(rec []string) bool {
	for _, v := range rec {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

// require checks every non-empty column name exists so a changed export format
// fails loudly instead of silently producing empty values.
func (t *table) require(columns ...string) error {
	var missing []string
	for _, c := range columns {
		if c == "" {
			continue
		}
		if _, ok := t.header[normaliseHeader(c)]; !ok {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		have := make([]string, 0, len(t.header))
		for h := range t.header {
			have = append(have, h)
		}
		return fmt.Errorf("%s: missing columns %q (found %q)", t.path, missing, have)
	}
	return nil
}

func (t *table) value(row []string, column string) string {
	if column == "" {
		return ""
	}
	i, ok := t.header[normaliseHeader(column)]
	if !ok || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}

var moneyJunk = regexp.MustCompile(`[£$€,\s]|GBP|USD|EUR`)

// parseMoney handles "£1,234.56", "-12.00", "(12.00)" and "12.00 DR"/"12.00 CR".
func parseMoney(s string) (decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "--" {
		return decimal.Zero, nil
	}
	neg := false
	upper := strings.ToUpper(s)
	switch {
	case strings.HasSuffix(upper, "DR"):
		neg = true
		s = s[:len(s)-2]
	case strings.HasSuffix(upper, "CR"):
		s = s[:len(s)-2]
	}
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		neg = true
		s = s[1 : len(s)-1]
	}
	s = moneyJunk.ReplaceAllString(s, "")
	s = strings.TrimPrefix(s, "+")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return d, fmt.Errorf("cannot parse amount %q", s)
	}
	if neg {
		d = d.Neg()
	}
	return d, nil
}

func parseDate(s string, formats []string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if len(formats) == 0 {
		formats = defaultDateFormats
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date %q (tried %q)", s, formats)
}
