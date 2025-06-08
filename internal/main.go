package internal

import "time"

// Transaction represents a single financial transaction, such as an entry from
// a bank statement.
type Transaction struct {
	// Date is the date and time when the transaction occurred.
	Date time.Time

	// NOTE: Consider adding BookingDate and ValueDate in stead

	// CounterpartName is the name of the person or entity receiving or sending funds.
	CounterpartName string

	// Description is an optional note or description providing additional details about the transaction.
	Description string

	// Amount is the signed amount of the transaction.
	// The value  is an integer that represents tha smalles currency unit (e.g., cents).
	Amount int

	// Currency is the ISO4217 code of the currency
	// Currency string
}
