package txn

// layout describes the mapping of transaction field names to their column positions
// in a specific format.
//
// Each canonical field name is associated with the index of
// the column where its value can be found. The index starts at 1. An index of 0
// indicates that the field is not present in the bank statement.
type layout struct {
	date    uint32
	payee   uint32
	memo    uint32
	amount  uint32
	inflow  uint32
	outflow uint32
}

type transaction struct {
	date    string
	payee   string
	memo    string
	inflow  string
	outflow string
}
