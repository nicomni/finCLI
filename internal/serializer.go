package txn

type serializer struct {
	lo layout
}

func (s serializer) serialize(tx transaction) []string {
	maxIdx := 5 // default to 5 fields
	indeces := []uint32{s.lo.date, s.lo.payee, s.lo.memo, s.lo.inflow, s.lo.outflow}
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
		{tx.date, s.lo.date},
		{tx.payee, s.lo.payee},
		{tx.memo, s.lo.memo},
		{tx.inflow, s.lo.inflow},
		{tx.outflow, s.lo.outflow},
	}
	for _, f := range fields {
		if f.index > 0 {
			out[f.index-1] = f.val
		}
	}
	return out
}
