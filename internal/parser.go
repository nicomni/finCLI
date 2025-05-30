package txn

import (
	"fmt"
)

type parser struct {
	lo Layout
}

func (p parser) parse(record []string) (Transaction, []error) {
	var errors []error
	var tx Transaction

	fields := []struct {
		name  string
		index uint32
		set   func(string)
	}{
		{"date", p.lo.Date, func(v string) { tx.Date = v }},
		{"payee", p.lo.Payee, func(v string) { tx.Payee = v }},
		{"memo", p.lo.Memo, func(v string) { tx.Memo = v }},
		{"inflow", p.lo.Inflow, func(v string) { tx.Inflow = v }},
		{"outflow", p.lo.Outflow, func(v string) { tx.Outflow = v }},
	}

	for _, f := range fields {
		if f.index == 0 {
			f.set("")
			continue
		}
		if f.index > uint32(len(record)) {
			errors = append(errors, fmt.Errorf("layout index %d for field %s out of range for record of length %d", f.index, f.name, len(record)))
			f.set("")
			continue
		}
		f.set(record[f.index-1])

	}

	return tx, errors
}
