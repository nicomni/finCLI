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
			format.ColumnMappings,
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
	txn StatementTransaction,
	colmap []TransactionColumn,
) error {
	record := make([]string, len(colmap))
	for _, col := range colmap {
		var value string
		switch col.Kind {
		case FieldDate:
			value = txn.Date.Format(col.DateFormat)
		case FieldPayee:
			value = txn.Payee
		case FieldMemo:
			value = txn.Memo
		case FieldInflow:
			value = formatAmount(txn.Inflow)
		case FieldOutflow:
			value = formatAmount(txn.Outflow)
		default:
			panic("unknown format. should never happen")
		}
		record[col.Pos-1] = value
	}
	return writer.Write(record)
}

func formatAmount(value int) string {
	major := value / 100
	minor := value % 100
	return fmt.Sprintf("%d.%02d", major, minor)
}
