package bankcsv

import (
	"fmt"
	"time"
)

// Format describes the strucutre of a CSV statement file.
//
// As returned by [NewFormat].
type Format struct {
	Delimiter        rune
	DateFormat       string
	DecimalSeparator rune
	ColumnMappings   []TransactionColumn

	// If set to true, reading and writing transaction records with this format
	// will start on the first line of the file.
	noHeader bool
}

// NewFormat creates a new Format with sensible defaults.
//
// The default format:
//   - uses a comma (',') as delimiter
//   - uses date format "2006-01-02" (time.DateOnly)
//   - uses a dot ('.') as decimal separator for amounts, e.g. "100.00"
//   - includes column headers
//   - does not include any column mappings.
func NewFormat() *Format {
	return &Format{
		Delimiter:        ',',
		DateFormat:       time.DateOnly,
		DecimalSeparator: '.',
	}
}

func (f *Format) HasHeader() bool {
	return !f.noHeader
}

func (f *Format) SetHasHeader(hasHeader bool) {
	f.noHeader = !hasHeader
}

type TransactionColumn struct {
	Name string
	Kind FieldKind
	Pos  int // Column position, starts at 1 (one). A 0 or negative value means not present.
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

type FormatRegistry map[string]Format

type Factory struct {
	InitRegistry *FormatRegistry
}

func NewRegistry(factory *Factory) *FormatRegistry {
	if factory != nil && factory.InitRegistry != nil {
		return factory.InitRegistry
	}
	return &defaultRegistry
}

var defaultRegistry = FormatRegistry{
	"bulder": {
		Delimiter: ';', DateFormat: time.DateOnly,
		DecimalSeparator: ',',
		ColumnMappings: []TransactionColumn{
			{Name: "Dato", Kind: FieldDate, Pos: 1},
			{Name: "Tekst", Kind: FieldMemo, Pos: 9},
			{Name: "Inn på konto", Kind: FieldInflow, Pos: 2},
			{Name: "Ut fra konto", Kind: FieldOutflow, Pos: 3},
		},
	},
	"ynab": {
		Delimiter: ',', DateFormat: time.DateOnly,
		DecimalSeparator: '.',
		ColumnMappings: []TransactionColumn{
			{Name: "Date", Kind: FieldDate, Pos: 1},
			{Name: "Payee", Kind: FieldPayee, Pos: 2},
			{Name: "Memo", Kind: FieldMemo, Pos: 3},
			{Name: "Inflow", Kind: FieldInflow, Pos: 4},
			{Name: "Outflow", Kind: FieldOutflow, Pos: 5},
		},
	},
}

func (r FormatRegistry) Get(name string) (Format, error) {
	format, ok := r[name]
	if !ok {
		return Format{}, fmt.Errorf("format '%s' is unknown", name)
	}
	return format, nil
}
