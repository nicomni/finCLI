package txn

import (
	"reflect"
	"testing"
)

func Test_serializer_serialize_zeroLayout(t *testing.T) {
	ser := serializer{lo: layout{}}
	tx := transaction{
		date:    "2025-01-01",
		payee:   "Bob's Store",
		memo:    "stuff",
		inflow:  "100.00",
		outflow: "50.00",
	}
	got := ser.serialize(tx)
	want := []string{"", "", "", "", ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serializer.serialize() with zero-layout = %v, want %v", got, want)
	}
}

func Test_serializer_serialize_customLayout(t *testing.T) {
	ser := serializer{lo: layout{
		date:    2,
		payee:   4,
		memo:    1,
		inflow:  3,
		outflow: 5,
	}}
	tx := transaction{
		date:    "2025-01-01",
		payee:   "Bob's Store",
		memo:    "stuff",
		inflow:  "100.00",
		outflow: "50.00",
	}
	got := ser.serialize(tx)
	want := []string{"stuff", "2025-01-01", "100.00", "Bob's Store", "50.00"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serializer.serialize() with custom layout = %v, want %v", got, want)
	}
}

func Test_serializer_serialize_layoutWithLargeIndexes(t *testing.T) {
	ser := serializer{lo: layout{
		date:    2,
		payee:   6,
		memo:    1,
		inflow:  4,
		outflow: 8,
	}}
	tx := transaction{
		date:    "2025-01-01",
		payee:   "Bob's Store",
		memo:    "stuff",
		inflow:  "100.00",
		outflow: "50.00",
	}
	got := ser.serialize(tx)
	want := []string{"stuff", "2025-01-01", "", "100.00", "", "Bob's Store", "", "50.00"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serializer.serialize() with large indexes = %v, want %v", got, want)
	}
}
