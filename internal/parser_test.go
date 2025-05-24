package txn

import (
	"regexp"
	"testing"
)

func Test_parser(t *testing.T) {
	tests := []struct {
		name   string
		layout layout
		record []string
		want   transaction
	}{
		{
			name:   "zero-layout",
			layout: layout{},
			record: []string{"2025-01-01", "Bob's Store", "stuff", "100.00", ""},
			want:   transaction{},
		},
		{
			name: "layout 1",
			layout: layout{
				date:    1,
				payee:   2,
				memo:    3,
				inflow:  4,
				outflow: 5,
			},
			record: []string{"2025-01-01", "Bob's Store", "stuff", "100.00", ""},
			want: transaction{
				date:    "2025-01-01",
				payee:   "Bob's Store",
				memo:    "stuff",
				inflow:  "100.00",
				outflow: "",
			},
		},
		{
			name: "layout 2",
			layout: layout{
				date:    2,
				payee:   4,
				memo:    1,
				inflow:  3,
				outflow: 5,
			},
			record: []string{"stuff", "2025-01-01", "100.00", "Bob's Store", ""},
			want: transaction{
				date:    "2025-01-01",
				payee:   "Bob's Store",
				memo:    "stuff",
				inflow:  "100.00",
				outflow: "",
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
	lo := layout{
		date:    1,
		payee:   2,
		memo:    6, // index larger than record length
		inflow:  3,
		outflow: 7, // index larger than record length
	}

	record := []string{"2025-01-01", "My Employer", "100.00"}
	want := transaction{
		date:    "2025-01-01",
		payee:   "My Employer",
		memo:    "",
		inflow:  "100.00",
		outflow: "",
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
	lo := layout{
		date:    3,
		payee:   5,
		memo:    1,
		inflow:  7,
		outflow: 2,
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
	want := transaction{
		date:    "2025-01-01",
		payee:   "Bob's Store",
		memo:    "foo",
		inflow:  "100.00",
		outflow: "bar",
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
