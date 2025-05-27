package txn

import (
	"fmt"
)

type parser struct {
	lo layout
}

func (p parser) parse(record []string) (transaction, []error) {
	var errors []error
	var tx transaction

	fields := []struct {
		name  string
		index uint32
		set   func(string)
	}{
		{"date", p.lo.Date, func(v string) { tx.date = v }},
		{"payee", p.lo.Payee, func(v string) { tx.payee = v }},
		{"memo", p.lo.Memo, func(v string) { tx.memo = v }},
		{"inflow", p.lo.Inflow, func(v string) { tx.inflow = v }},
		{"outflow", p.lo.Outflow, func(v string) { tx.outflow = v }},
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
