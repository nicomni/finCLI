package txn

type serializer struct {
	lo layout
}

func (s serializer) serialize(tx transaction) []string {
	maxIdx := 5 // default to 5 fields
	indeces := []uint32{s.lo.Date, s.lo.Payee, s.lo.Memo, s.lo.Inflow, s.lo.Outflow}
	for _, idx := range indeces {
		if int(idx) > maxIdx {
			maxIdx = int(idx)
		}
	}
	out := make([]string, maxIdx)

	fields := []struct {
		val   string
		index uint32
	}{
		{tx.date, s.lo.Date},
		{tx.payee, s.lo.Payee},
		{tx.memo, s.lo.Memo},
		{tx.inflow, s.lo.Inflow},
		{tx.outflow, s.lo.Outflow},
	}
	for _, f := range fields {
		if f.index > 0 {
			out[f.index-1] = f.val
		}
	}
	return out
}
