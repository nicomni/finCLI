package csvstatement

import (
	"fmt"
	"time"
)

// Format describes the strucutre of a CSV statement file.
//
// For convenience use [NewFormat] which sets sensible defaults.
type Format struct {
	Id             string
	Delimiter      rune
	HasHeader      bool
	ColumnMappings []TransactionColumn
}

// NewFormat returns a Format with HasHeader set to true by default.
func NewFormat() Format {
	return Format{HasHeader: true}
}

type TransactionColumn struct {
	Name       string
	Kind       FieldKind
	Pos        int // Column position, starts at 1 (one). A 0 or negative value means not present.
	DateFormat string
}

// FieldKind describes the kind of data in a column
type FieldKind string

const (
	FieldDate    FieldKind = "date"
	FieldPayee   FieldKind = "payee"
	FieldMemo    FieldKind = "memo"
	FieldInflow  FieldKind = "inflow"
	FieldOutflow FieldKind = "outflow"
)

var formatRegistry = map[string]Format{
	"bulder": {
		Id: "bulder", Delimiter: ';', HasHeader: true,
		ColumnMappings: []TransactionColumn{
			{Name: "Dato", Kind: FieldDate, Pos: 1, DateFormat: time.DateOnly},
			{Name: "Tekst", Kind: FieldMemo, Pos: 9},
			{Name: "Inn på konto", Kind: FieldInflow, Pos: 2},
			{Name: "Ut fra konto", Kind: FieldOutflow, Pos: 3},
		},
	},
	"ynab": {
		Id: "ynab", Delimiter: ',', HasHeader: true,
		ColumnMappings: []TransactionColumn{
			{Name: "Date", Kind: FieldDate, Pos: 1, DateFormat: time.DateOnly},
			{Name: "Payee", Kind: FieldPayee, Pos: 2},
			{Name: "Memo", Kind: FieldMemo, Pos: 3},
			{Name: "Inflow", Kind: FieldInflow, Pos: 4},
			{Name: "Outflow", Kind: FieldOutflow, Pos: 5},
		},
	},
}

func GetFormat(name string) (Format, error) {
	format, ok := formatRegistry[name]
	if !ok {
		return Format{}, fmt.Errorf("format '%s' is unknown", name)
	}
	return format, nil
}
