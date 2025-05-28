package txn

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// layout describes the mapping of transaction field names to their column positions
// in a specific format.
//
// Each canonical field name is associated with the index of
// the column where its value can be found. The index starts at 1. An index of 0
// indicates that the field is not present in the bank statement.
type layout struct {
	Date    uint32
	Payee   uint32
	Memo    uint32
	Inflow  uint32
	Outflow uint32
}

type transaction struct {
	date    string
	payee   string
	memo    string
	inflow  string
	outflow string
}

//go:embed layouts/*
var layouts embed.FS

func getLayout(name string) (layout, error) {
	layoutBytes, err := layouts.ReadFile(filepath.Join("layouts", name+".yaml"))
	if err != nil {
		return layout{}, fmt.Errorf("could not read builtin layout: %v", err)
	}
	result, err := unmarshalLayout(layoutBytes)
	if err != nil {
		return result, fmt.Errorf("could not decode builtin layout: %v", err)
	}
	return result, nil
}

func unmarshalLayout(layoutBytes []byte) (layout, error) {
	var result layout
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewReader(layoutBytes))
	if err != nil {
		return layout{}, err
	}
	// NOTE: Consider better error handling for unknown fields if strict decoding is required
	err = v.UnmarshalExact(&result)
	return result, err
}
