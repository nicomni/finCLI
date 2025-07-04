package bankcsv

import (
	"encoding/csv"
	"errors"
	"fincli/internal/domain"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type CSVStatement struct {
	Transactions []domain.Transaction

	format *Format
}

type Parser struct {
	log    slog.Logger
	format *Format
}

func NewParser(format *Format) *Parser {
	var parser Parser
	parser.log = *slog.Default()

	if len(format.ColumnMappings) == 0 {
		parser.log.Warn("Creating parser with empty column mapping. Parsing will return empty transactions")
	}
	parser.format = format
	return &parser
}

func (p Parser) Parse(source io.Reader) (*CSVStatement, error) {
	// TODO: Validate that input conforms to format, and is not empty.
	result := new(CSVStatement)

	reader := csv.NewReader(source)
	if p.format.Delimiter != 0 {
		reader.Comma = p.format.Delimiter
	}

	if p.format.HasHeader() {
		// TEST: Without header
		_, err := reader.Read()
		if err != nil {
			return nil, fmt.Errorf("parsing statement: could not read header. Error: %w", err)
		}
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read records. Error: %w", err)
	}

	p.checkColumnMappings(reader.FieldsPerRecord)

	result.Transactions = make([]domain.Transaction, 0, len(records))

	for _, rec := range records {
		csvTxn, err := p.parseCsvRecord(rec)
		if err != nil {
			return nil, err
		}
		result.Transactions = append(result.Transactions, csvTxn.Transaction)
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

func (p Parser) parseCsvRecord(record []string) (*CSVTransaction, error) {
	var txn domain.Transaction
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
			txn.Amount += abs(amount)
		case FieldOutflow:
			amount, err := strconv.Atoi(normalizeDecimal(value))
			if err != nil {
				return nil, fmt.Errorf("could not parse outflow value at column position '%d' with value '%s': %w", col.Pos, value, err)
			}
			txn.Amount -= abs(amount)
		case FieldAmount:
			amount, err := strconv.Atoi(normalizeDecimal(value))
			if err != nil {
				return nil, fmt.Errorf("could not parse amount value at column position '%d' with value '%s': %w", col.Pos, value, err)
			}
			txn.Amount = amount
		default:
			return nil, fmt.Errorf("could not parse record: unknown field kind '%s' in column '%s'", col.Kind, col.Name)
		}
	}
	csvTxn := new(CSVTransaction)
	csvTxn.Transaction = txn
	csvTxn.Format = p.format

	return csvTxn, nil
}

// normalizeDecimal removes all spaces, commas, dots, and redundant sign
// character (+) from the input string.
//
// This is useful for normalizing decimal numbers (assuming precision 2)
// from various locales so they can be parsed as integers representing the
// smallest currency unit (e.g., cents).
func normalizeDecimal(dirty string) (clean string) {
	// HACK: Breaks the program if input does not have a precision of 2.
	clean = strings.ReplaceAll(dirty, " ", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.ReplaceAll(clean, "+", "")
	return clean
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
