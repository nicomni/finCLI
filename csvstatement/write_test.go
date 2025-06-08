package csvstatement_test

import (
	"fincli/csvstatement"
	"strings"
	"testing"
	"time"
)

func Test_WriteToYNAB(t *testing.T) {
	statement := csvstatement.ParsedStatement{
		Transactions: []csvstatement.StatementTransaction{
			{
				Date:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Payee:   "testPayee",
				Memo:    "testMemo",
				Inflow:  0,
				Outflow: 1234,
			},
			{
				Date:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				Payee:   "testPayee2",
				Memo:    "testMemo2",
				Inflow:  50000,
				Outflow: 0,
			},
		},
	}

	format, err := csvstatement.GetFormat("ynab")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `Date,Payee,Memo,Inflow,Outflow
2025-01-01,testPayee,testMemo,0.00,12.34
2025-01-02,testPayee2,testMemo2,500.00,0.00
`

	w := strings.Builder{}
	err = csvstatement.WriteStatement(&w, statement, format)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	result := w.String()
	if result != want {
		t.Errorf("Write statement, unexpected output:\ngot:\n%s\nbut wanted:\n%s", result, want)
	}
}
