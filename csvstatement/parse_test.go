package csvstatement_test

import (
	"strings"
	"testing"
	"time"
	"txn/csvstatement"
)

func TestParser_Basic(t *testing.T) {
	format := csvstatement.Format{
		Delimiter: ',',
		HasHeader: true,
		ColumnMappings: []csvstatement.TransactionColumn{
			{Name: "Date", Kind: csvstatement.FieldDate, Pos: 1, DateFormat: "2006-01-02"},
			{Name: "Payee", Kind: csvstatement.FieldPayee, Pos: 2},
			{Name: "Memo", Kind: csvstatement.FieldMemo, Pos: 3},
			{Name: "Inflow", Kind: csvstatement.FieldInflow, Pos: 4},
			{Name: "Outflow", Kind: csvstatement.FieldOutflow, Pos: 5},
		},
	}
	csvData := `Date, Payee, Memo, Inflow, Outflow
2025-01-01,Store,Groceries,0.00,12.34
2025-01-02,Bankomat,Deposit,500.00,0.00`

	parser := csvstatement.NewParser(format)

	result, err := parser.Parse(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(result.Transactions))
	}
	tx1 := result.Transactions[0]
	want := csvstatement.StatementTransaction{
		Date:    time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Payee:   "Store",
		Memo:    "Groceries",
		Outflow: 1234,
	}
	if tx1 != want {
		t.Errorf("unexpected tx1: %+v, wanted: %+v", tx1, want)
	}

	tx2 := result.Transactions[1]
	want = csvstatement.StatementTransaction{
		Date:   time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
		Payee:  "Bankomat",
		Memo:   "Deposit",
		Inflow: 50000,
	}
	if tx2 != want {
		t.Errorf("unexpected tx2: %+v, wanted: %+v", tx1, want)
	}
}

func Test_Bulder(t *testing.T) {
	csvData := `Dato;Inn på konto;Ut fra konto;Til konto;Til kontonummer;Fra konto;Fra kontonummer;Type;Tekst;KID;Hovedkategori;Underkategori
2025-01-01;;12,34;;;;;;Groceries;;;
2025-01-02;500,00;;;;;;;Deposit;;;`
	format, err := csvstatement.GetFormat("bulder")
	if err != nil {
		t.Fatal(err)
	}
	parser := csvstatement.NewParser(format)
	got, err := parser.Parse(strings.NewReader(csvData))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(got.Transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(got.Transactions))
	}
	tx1 := got.Transactions[0]
	want := csvstatement.StatementTransaction{
		Date:    time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Memo:    "Groceries",
		Outflow: 1234,
	}
	if tx1 != want {
		t.Errorf("unexpected tx1: %+v, wanted: %+v", tx1, want)
	}

	tx2 := got.Transactions[1]
	want = csvstatement.StatementTransaction{
		Date:   time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
		Memo:   "Deposit",
		Inflow: 50000,
	}
	if tx2 != want {
		t.Errorf("unexpected tx2: %+v, wanted: %+v", tx1, want)
	}
}
