package csvstatement_test

import (
	"fincli/internal/csvstatement"
	"fincli/internal/domain"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestParser_Basic(t *testing.T) {
	format := csvstatement.Format{
		Delimiter:        ',',
		HasHeader:        true,
		DateFormat:       time.DateOnly,
		DecimalSeparator: '.',
		ColumnMappings: []csvstatement.TransactionColumn{
			{Name: "Date", Kind: csvstatement.FieldDate, Pos: 1},
			{Name: "Payee", Kind: csvstatement.FieldPayee, Pos: 2},
			{Name: "Memo", Kind: csvstatement.FieldMemo, Pos: 3},
			{Name: "Inflow", Kind: csvstatement.FieldInflow, Pos: 4},
			{Name: "Outflow", Kind: csvstatement.FieldOutflow, Pos: 5},
		},
	}
	csvData := (func() string {
		b := strings.Builder{}
		b.WriteString("Date, Payee, Memo, Inflow, Outflow\n")
		b.WriteString("2025-01-01,Store,Groceries,0.00,12.34\n")
		b.WriteString("2025-01-02,Bankomat,Deposit,500.00,0.00")
		return b.String()
	})()

	wantTxns := []domain.Transaction{
		{
			Date:            time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			CounterpartName: "Store",
			Description:     "Groceries",
			Amount:          -1234,
		},
		{
			Date:            time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
			CounterpartName: "Bankomat",
			Description:     "Deposit",
			Amount:          50000,
		},
	}
	parser := csvstatement.NewParser(format)

	result, err := parser.Parse(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(result.Transactions))
	}
	for idx, txn := range result.Transactions {
		if err := checkEqual(wantTxns[idx], txn); err != nil {
			t.Errorf("Transaction %d: %v", idx, err)
		}
	}
}

func Test_Bulder(t *testing.T) {
	csvData := (func() string {
		b := strings.Builder{}
		b.WriteString("Dato;Inn på konto;Ut fra konto;Til konto;Til kontonummer;" +
			"Fra konto;Fra kontonummer;Type;Tekst;KID;Hovedkategori;Underkategori\n")
		b.WriteString("2025-01-01;;12,34;;;;;;Groceries;;;\n")
		b.WriteString("2025-01-02;500,00;;;;;;;Deposit;;;")
		return b.String()
	})()

	wantTxns := []domain.Transaction{
		{
			Date:            time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			CounterpartName: "",
			Description:     "Groceries",
			Amount:          -1234,
		},
		{
			Date:            time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
			CounterpartName: "",
			Description:     "Deposit",
			Amount:          50000,
		},
	}

	registry := csvstatement.NewRegistry(nil)
	format, err := registry.Get("bulder")
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

	for idx, txn := range got.Transactions {
		if err := checkEqual(wantTxns[idx], txn); err != nil {
			t.Errorf("Transaction %d: %v", idx, err)
		}
	}
}

func checkEqual(want, got domain.Transaction) error {
	var errs []string

	if got.Date != want.Date {
		errs = append(errs, fmt.Sprintf("Date field: want %v, got %v", want.Date, got.Date))
	}

	if got.CounterpartName != want.CounterpartName {
		errs = append(errs, fmt.Sprintf("Counterpart name: want %v, got %v", want.CounterpartName, got.CounterpartName))
	}

	if got.Description != want.Description {
		errs = append(errs, fmt.Sprintf("Description field: want %v, got %v", want.Description, got.Description))
	}

	if got.Amount != want.Amount {
		errs = append(errs, fmt.Sprintf("Amount field: want %d, got %d", want.Amount, got.Amount))
	}

	if len(errs) > 0 {
		return fmt.Errorf("Unexpected transaction field values:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}
