package csvstatement

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type ParsedStatement struct {
	Transactions []StatementTransaction
}

type Parser struct {
	log    slog.Logger
	format Format
}

func NewParser(format Format) *Parser {
	var parser Parser
	parser.log = *slog.Default()

	if len(format.ColumnMappings) == 0 {
		parser.log.Warn("Creating parser with empty column mapping. Parsing will return empty transactions")
	}
	parser.format = format
	return &parser
}

func (p Parser) Parse(source io.Reader) (ParsedStatement, error) {
	var result ParsedStatement

	reader := csv.NewReader(source)
	if p.format.Delimiter != 0 {
		reader.Comma = p.format.Delimiter
	}

	if p.format.HasHeader {
		// TEST: Without header
		_, err := reader.Read()
		if err != nil {
			return ParsedStatement{}, fmt.Errorf("parsing statement: could not read header. Error: %w", err)
		}
	}

	records, err := reader.ReadAll()
	if err != nil {
		return result, fmt.Errorf("could not read records. Error: %w", err)
	}

	p.checkColumnMappings(reader.FieldsPerRecord)

	result.Transactions = make([]StatementTransaction, 0, len(records))

	for _, rec := range records {
		txn, err := p.parseCsvRecord(rec)
		if err != nil {
			return ParsedStatement{}, err
		}
		result.Transactions = append(result.Transactions, *txn)
	}

	return result, nil
}

// ErrNoColumnMap is returned when the format is not properly configured.
var ErrNoColumnMap = errors.New("format has no column mappings")

// checkColumnMappings validates the configured column positions in the Format
// against the actual number of fields in a CSV record.
//
// If a column mapping specifies a position greater than the number of
// available fields, a warning is logged. If a column mapping specifies a
// position of 0 or less, an info message is logged indicating the field will
// be skipped.
func (p *Parser) checkColumnMappings(numOfFields int) {
	p.log.Info("Validating column mapping against CSV")
	for _, col := range p.format.ColumnMappings {
		if col.Pos > numOfFields {
			p.log.Warn(
				fmt.Sprintf("Warning: Column '%s', with field type '%s', "+
					"is mapped to position %d, but the CSV record only has %d fields. "+
					"This column will be ignored.",
					col.Name, col.Kind, col.Pos, numOfFields,
				))
		}

		if col.Pos <= 0 {
			p.log.Info(fmt.Sprintf(
				"Column '%s' has position %d and will be skipped.",
				col.Name, col.Pos,
			))
		}
	}
}

func (p Parser) parseCsvRecord(record []string) (*StatementTransaction, error) {
	var txn StatementTransaction
	colMap := p.format.ColumnMappings
	for _, col := range colMap {
		if col.Pos <= 0 {
			continue
		}
		if col.Pos > len(record) {
			continue
		}
		value := record[col.Pos-1]
		if value == "" {
			continue
		}
		switch col.Kind {
		case FieldDate:
			date, err := time.Parse(p.format.DateFormat, value)
			if err != nil {
				return nil, fmt.Errorf("could not parse date in colum '%s' at position '%d': value='%s'. Error: %w", col.Name, col.Pos, value, err)
			}
			txn.Date = date
		case FieldPayee:
			txn.CounterpartName = value
		case FieldMemo:
			txn.Description = value
		case FieldInflow:
			amount, err := strconv.Atoi(normalizeDecimal(value))
			if err != nil {
				return nil, fmt.Errorf("could not parse inflow value at column position '%d' with value '%s': %w", col.Pos, value, err)
			}
			txn.Amount += amount
		case FieldOutflow:
			amount, err := strconv.Atoi(normalizeDecimal(value))
			if err != nil {
				return nil, fmt.Errorf("could not parse outflow value at column position '%d' with value '%s': %w", col.Pos, value, err)
			}
			txn.Amount -= amount
		}
	}

	return &txn, nil
}

// normalizeDecimal removes all spaces, commas, dots, and sign characters (+/-)
// from the input string.
//
// This is useful for normalizing decimal numbers (with assumed precision 2)
// from various locales so they can be parsed as integers representing the
// smallest currency unit (e.g., cents).
func normalizeDecimal(dirty string) (clean string) {
	clean = strings.ReplaceAll(dirty, " ", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ReplaceAll(clean, "+", "")
	clean = strings.ReplaceAll(clean, "-", "")
	return clean
}
