package csvstatement_test

import (
	"fincli/internal/csvstatement"
	"fincli/internal/domain"
	"strings"
	"testing"
	"time"
)

func Test_Write(t *testing.T) {
	statement := csvstatement.ParsedStatement{
		Transactions: []domain.Transaction{
			{
				Date:            time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				CounterpartName: "testPayee",
				Description:     "testMemo",
				Amount:          -1234,
			},
			{
				Date:            time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				CounterpartName: "testPayee2",
				Description:     "testMemo2",
				Amount:          50000,
			},
		},
	}

	tests := []struct {
		name     string
		formatId string
		want     string
	}{
		{
			name:     "YNAB",
			formatId: "ynab",
			want: (func() string {
				b := strings.Builder{}
				b.WriteString("Date,Payee,Memo,Inflow,Outflow\n")
				b.WriteString("2025-01-01,testPayee,testMemo,0.00,12.34\n")
				b.WriteString("2025-01-02,testPayee2,testMemo2,500.00,0.00\n")
				return b.String()
			})(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := csvstatement.NewRegistry(nil)
			format, err := registry.Get(tt.formatId)
			if err != nil {
				t.Fatal(err)
			}

			writer := strings.Builder{}
			err = csvstatement.WriteStatement(&writer, statement, format)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			result := writer.String()

			if result != tt.want {
				t.Errorf("Unexpected output. Most likely explained by following messages.")
			}

			gotLines := strings.Split(result, "\n")

			if len(gotLines) != 4 {
				t.Fatalf("Expected output to have 4 lines, but it had %d", len(gotLines))
			}

			wantLines := strings.Split(tt.want, "\n")

			for i := range gotLines {
				if gotLines[i] != wantLines[i] {
					t.Errorf("Line %d mismatch:\ngot:\t%q\nwant\t%q", i+1, gotLines[i], wantLines[i])
				}
			}

			if !strings.HasSuffix(result, "\n") {
				t.Errorf("Output does not end with a newline")
			}
		})
	}
}
