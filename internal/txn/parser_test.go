package txn

import (
	"regexp"
	"testing"
)

func Test_parser(t *testing.T) {
	tests := []struct {
		name   string
		layout Layout
		record []string
		want   Transaction
	}{
		{
			name:   "zero-layout",
			layout: Layout{},
			record: []string{"2025-01-01", "Bob's Store", "stuff", "100.00", ""},
			want:   Transaction{},
		},
		{
			name: "layout 1",
			layout: Layout{
				Date:    1,
				Payee:   2,
				Memo:    3,
				Inflow:  4,
				Outflow: 5,
			},
			record: []string{"2025-01-01", "Bob's Store", "stuff", "100.00", ""},
			want: Transaction{
				Date:    "2025-01-01",
				Payee:   "Bob's Store",
				Memo:    "stuff",
				Inflow:  "100.00",
				Outflow: "",
			},
		},
		{
			name: "layout 2",
			layout: Layout{
				Date:    2,
				Payee:   4,
				Memo:    1,
				Inflow:  3,
				Outflow: 5,
			},
			record: []string{"stuff", "2025-01-01", "100.00", "Bob's Store", ""},
			want: Transaction{
				Date:    "2025-01-01",
				Payee:   "Bob's Store",
				Memo:    "stuff",
				Inflow:  "100.00",
				Outflow: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser{tt.layout}
			got, errs := p.parse(tt.record)
			if len(errs) > 0 {
				t.Fatalf("parseTransaction() failed: %v", errs)
			}
			if got != tt.want {
				t.Errorf("parseTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_layoutWithLargeIndexes(t *testing.T) {
	lo := Layout{
		Date:    1,
		Payee:   2,
		Memo:    6, // index larger than record length
		Inflow:  3,
		Outflow: 7, // index larger than record length
	}

	record := []string{"2025-01-01", "My Employer", "100.00"}
	want := Transaction{
		Date:    "2025-01-01",
		Payee:   "My Employer",
		Memo:    "",
		Inflow:  "100.00",
		Outflow: "",
	}
	p := parser{lo: lo}
	got, errs := p.parse(record)

	t.Run("parses available fields", func(t *testing.T) {
		if got != want {
			t.Errorf("parse() = %v, wnat %v", got, want)
		}
	})

	t.Run("returns errors for out-of-range indexes", func(t *testing.T) {
		if len(errs) != 2 {
			t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
		}
	})

	t.Run("error messages are correct", func(t *testing.T) {
		re := regexp.MustCompile(`^layout index \d+ for field \w+ out of range for record of length \d+$`)
		for _, err := range errs {
			if err == nil || !re.MatchString(err.Error()) {
				t.Errorf("unexpected error message: %v", err)
			}
		}
	})
}

func Test_parser_parsesFromLargeRecord(t *testing.T) {
	lo := Layout{
		Date:    3,
		Payee:   5,
		Memo:    1,
		Inflow:  7,
		Outflow: 2,
	}
	record := []string{
		"foo",         // 1: memo
		"bar",         // 2: outflow
		"2025-01-01",  // 3: date
		"baz",         // 4: ignored
		"Bob's Store", // 5: payee
		"qux",         // 6: ignored
		"100.00",      // 7: inflow
		"extra",       // 8: ignored
	}
	want := Transaction{
		Date:    "2025-01-01",
		Payee:   "Bob's Store",
		Memo:    "foo",
		Inflow:  "100.00",
		Outflow: "bar",
	}
	p := parser{lo}
	got, errs := p.parse(record)
	if len(errs) > 0 {
		t.Fatalf("parse() returned errors: %v", errs)
	}
	if got != want {
		t.Errorf("parse() = %v, want %v", got, want)
	}
}
