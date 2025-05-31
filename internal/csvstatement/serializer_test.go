package csvstatement

import (
	"reflect"
	"testing"
)

func Test_serializer_serialize_zeroLayout(t *testing.T) {
	ser := serializer{lo: Layout{}}
	tx := Transaction{
		Date:    "2025-01-01",
		Payee:   "Bob's Store",
		Memo:    "stuff",
		Inflow:  "100.00",
		Outflow: "50.00",
	}
	got := ser.serialize(tx)
	want := []string{"", "", "", "", ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serializer.serialize() with zero-layout = %v, want %v", got, want)
	}
}

func Test_serializer_serialize_customLayout(t *testing.T) {
	ser := serializer{lo: Layout{
		Date:    2,
		Payee:   4,
		Memo:    1,
		Inflow:  3,
		Outflow: 5,
	}}
	tx := Transaction{
		Date:    "2025-01-01",
		Payee:   "Bob's Store",
		Memo:    "stuff",
		Inflow:  "100.00",
		Outflow: "50.00",
	}
	got := ser.serialize(tx)
	want := []string{"stuff", "2025-01-01", "100.00", "Bob's Store", "50.00"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serializer.serialize() with custom layout = %v, want %v", got, want)
	}
}

func Test_serializer_serialize_layoutWithLargeIndexes(t *testing.T) {
	ser := serializer{lo: Layout{
		Date:    2,
		Payee:   6,
		Memo:    1,
		Inflow:  4,
		Outflow: 8,
	}}
	tx := Transaction{
		Date:    "2025-01-01",
		Payee:   "Bob's Store",
		Memo:    "stuff",
		Inflow:  "100.00",
		Outflow: "50.00",
	}
	got := ser.serialize(tx)
	want := []string{"stuff", "2025-01-01", "", "100.00", "", "Bob's Store", "", "50.00"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serializer.serialize() with large indexes = %v, want %v", got, want)
	}
}
