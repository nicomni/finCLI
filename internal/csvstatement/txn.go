package csvstatement

import (
	"bytes"
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/spf13/viper"
)

// Layout describes the mapping of transaction field names to their column positions
// in a specific format.
//
// Each canonical field name is associated with the index of
// the column where its value can be found. The index starts at 1. An index of 0
// indicates that the field is not present in the bank statement.
type Layout struct {
	Date    uint32
	Payee   uint32
	Memo    uint32
	Inflow  uint32
	Outflow uint32
}

type Transaction struct {
	Date    string
	Payee   string
	Memo    string
	Inflow  string
	Outflow string
}

func LoadStatement(r io.Reader, layoutId string) ([]Transaction, error) {
	l, err := GetLayout(layoutId)
	if err != nil {
		return nil, err
	}
	p := parser{lo: l}

	cr := csv.NewReader(r)
	cr.Comma = ';'
	records, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read bank statement records: %v", err)
	}

	transactions := make([]Transaction, 0, len(records))

	for idx, rec := range records {
		if idx == 0 {
			continue // skip header
		}
		t, errs := p.parse(rec)
		// TODO: Lenient but loud: print warnings and continue
		if len(errs) != 0 {
			return nil, fmt.Errorf("parsing record %d: got errors: %v", idx, errs)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

//go:embed layouts/*
var layouts embed.FS

func GetLayout(name string) (Layout, error) {
	layoutBytes, err := layouts.ReadFile(filepath.Join("layouts", name+".yaml"))
	if errors.Is(err, fs.ErrNotExist) {
		return Layout{}, fmt.Errorf("layout '%s' is unknown", name)
	}
	if err != nil {
		return Layout{}, fmt.Errorf("could not read builtin layout: %v", err)
	}
	result, err := unmarshalLayout(layoutBytes)
	if err != nil {
		return result, fmt.Errorf("could not decode builtin layout: %v", err)
	}
	return result, nil
}

func unmarshalLayout(layoutBytes []byte) (Layout, error) {
	var result Layout
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewReader(layoutBytes))
	if err != nil {
		return Layout{}, err
	}
	// NOTE: Consider better error handling for unknown fields if strict decoding is required
	err = v.UnmarshalExact(&result)
	return result, err
}
