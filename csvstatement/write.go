package csvstatement

import (
	"encoding/csv"
	"fmt"
	"io"
)

func WriteStatement(writer io.Writer, statement ParsedStatement, format Format) error {
	csvwriter := csv.NewWriter(writer)
	defer csvwriter.Flush()
	if format.HasHeader {
		if err := writeHeader(csvwriter, format.ColumnMappings); err != nil {
			return fmt.Errorf("could not write CSV header: %w", err)
		}
	}

	for idx, txn := range statement.Transactions {
		if err := writeRecord(
			csvwriter,
			txn,
			format,
		); err != nil {
			return fmt.Errorf("could not write transaction %d as CSV record: %w", idx, err)
		}
	}
	return nil
}

func writeHeader(writer *csv.Writer, colmap []TransactionColumn) error {
	headers := make([]string, len(colmap))
	for _, col := range colmap {
		headers[col.Pos-1] = col.Name
	}
	return writer.Write(headers)
}

func writeRecord(
	writer *csv.Writer,
	txn Transaction,
	format Format,
) error {
	record, err := constructRecord(txn, format)
	if err != nil {
		return err
	}
	return writer.Write(record)
}

func constructRecord(txn Transaction, format Format) ([]string, error) {
	colMap := format.ColumnMappings
	record := make([]string, len(colMap))
	for _, col := range colMap {
		var value string
		switch col.Kind {
		case FieldDate:
			value = txn.Date.Format(format.DateFormat)
		case FieldPayee:
			value = txn.CounterpartName
		case FieldMemo:
			value = txn.Description
		case FieldInflow:
			if txn.Amount > 0 {
				value = formatAmount(txn.Amount, format)
			} else {
				value = formatAmount(0, format)
			}
		case FieldOutflow:
			if txn.Amount < 0 {
				value = formatAmount(-txn.Amount, format)
			} else {
				value = formatAmount(0, format)
			}
		default:
			return nil, fmt.Errorf("could not construct record field: unknown field kind '%s'", col.Kind)
		}
		record[col.Pos-1] = value
	}
	return record, nil
}

func formatAmount(value int, format Format) string {
	major := value / 100
	minor := value % 100
	return fmt.Sprintf("%d%c%02d", major, format.DecimalSeparator, minor)
}
