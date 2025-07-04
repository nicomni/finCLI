// Package bankcsv provides tools for reading, writing, formatting, bank
// statements in CSV file formats.
package bankcsv

import "fincli/internal/domain"

// TODO: 1. We should be able to access fields by header name, if the CSV format has headers.
// TODO: 2. We should be able to access fields by index of the underlying CSV record.

// CSVTransaction represents the run-time representation of a parsed CSV
// transaction.
type CSVTransaction struct {
	Transaction domain.Transaction // The parsed transaction data

	Format *Format // The format used to parse the transaction
}
