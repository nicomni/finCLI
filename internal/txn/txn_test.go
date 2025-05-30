package txn_test

import (
	"bytes"
	"reflect"
	"testing"
	"txn/internal/txn"
)

func TestLoadStatement(t *testing.T) {
	t.Run("bulder", func(t *testing.T) {
		stmt := []byte(`Dato;Inn på konto;Ut fra konto;Til konto;Til kontonummer;Fra konto;Fra kontonummer;Type;Tekst;KID;Hovedkategori;Underkategori
2025-01-01;123,45;;;;;;;ABC;;;
2025-01-02;;-234,56;;;;;;XYZ;;;`)
		want := []txn.Transaction{
			{Date: "2025-01-01", Payee: "", Memo: "ABC", Inflow: "123,45", Outflow: ""},
			{Date: "2025-01-02", Payee: "", Memo: "XYZ", Inflow: "", Outflow: "-234,56"},
		}
		reader := bytes.NewBuffer(stmt)
		got, err := txn.LoadStatement(reader, "bulder")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("wanted %v,\nbut got: %v", want, got)
		}
	})
}

func TestGetLayout(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		layoutName := "this_layout_does_not_exist"
		expectedLayout := txn.Layout{}
		expectedErrMsg := "layout 'this_layout_does_not_exist' is unknown"
		l, err := txn.GetLayout(layoutName)

		if err == nil {
			t.Error("expected error, but got nil")
		}
		if err.Error() != expectedErrMsg {
			t.Errorf("unexpected error message: %v", err)
		}
		if l != expectedLayout {
			t.Errorf("expected zero-layout, got %v", l)
		}
	})

	t.Run("bulder", func(t *testing.T) {
		expectedLayout := txn.Layout{Date: 1, Memo: 9, Payee: 0, Inflow: 2, Outflow: 3}

		got, err := txn.GetLayout("bulder")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != expectedLayout {
			t.Errorf("expected: %v, got: %v", expectedLayout, got)
		}
	})
}
