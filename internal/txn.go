package txn

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
