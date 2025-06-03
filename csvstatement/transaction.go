package csvstatement

import "time"

// Transaction represents a single financial transaction, such as an entry from
// a bank statement.
//
// It contains information about the date, payee, memo, and the amounts
// involved. Amounts are stored as int64 in the smallest currency unit (e.g.,
// cents) to ensure precision.
type StatementTransaction struct {
	// Date is the date and time when the transaction occurred.
	Date time.Time

	// Payee is the name of the person or entity receiving or sending funds.
	Payee string

	// Memo is an optional note or description providing additional details about the transaction.
	Memo string

	// Inflow is the amount of money received, in the smallest currency unit (e.g., cents).
	Inflow int

	// Outflow is the amount of money spent or withdrawn, in the smallest currency unit (e.g., cents).
	Outflow int
}
