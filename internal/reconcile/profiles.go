package reconcile

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Profiles is the file describing every CSV export format we know how to read.
// Adding a new bank or marketplace is a config change, not a code change.
type Profiles struct {
	Statements map[string]StatementProfile `json:"statements"`
	Purchases  map[string]PurchaseProfile  `json:"purchases"`
}

func LoadProfiles(path string) (*Profiles, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Profiles
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	for name, s := range p.Statements {
		s.Name = name
		p.Statements[name] = s
	}
	for name, s := range p.Purchases {
		s.Name = name
		if s.Source == "" {
			s.Source = name
		}
		p.Purchases[name] = s
	}
	return &p, nil
}

func (p *Profiles) Statement(name string) (StatementProfile, error) {
	s, ok := p.Statements[name]
	if !ok {
		return s, fmt.Errorf("unknown statement profile %q (have %q)", name, keys(p.Statements))
	}
	return s, nil
}

func (p *Profiles) Purchase(name string) (PurchaseProfile, error) {
	s, ok := p.Purchases[name]
	if !ok {
		return s, fmt.Errorf("unknown purchase profile %q (have %q)", name, keys(p.Purchases))
	}
	return s, nil
}

func keys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
