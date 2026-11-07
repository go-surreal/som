//go:build embed

package lib

// RecordRange restricts a query to a range of record IDs. From and To hold the
// rendered record ID keys of the bounds (e.g. "$A" or "[$A, $B]"); an empty
// key leaves that side of the range open.
type RecordRange struct {
	From          string
	FromInclusive bool
	To            string
	ToInclusive   bool
}

func (r RecordRange) render() string {
	return r.renderFrom() + ".." + r.renderTo()
}

func (r RecordRange) renderFrom() string {
	if r.From == "" || r.FromInclusive {
		return r.From
	}

	return r.From + ">"
}

func (r RecordRange) renderTo() string {
	if r.To == "" || !r.ToInclusive {
		return r.To
	}

	return "=" + r.To
}

// renderFromGuard returns a filter enforcing the lower bound of the range, or
// an empty string if the range is open at the start.
func (r RecordRange) renderFromGuard(table string) string {
	if r.From == "" {
		return ""
	}

	operator := OpGreaterThan
	if r.FromInclusive {
		operator = OpGreaterThanEqual
	}

	return "id " + string(operator) + " " + recordFromKey(table, r.From)
}

// seekAfter returns the range continuing exclusively after the record ID held
// by the given query variable, keeping the upper bound of the range.
func (r RecordRange) seekAfter(idVar string) RecordRange {
	return RecordRange{
		From:        "record::id(" + idVar + ")",
		To:          r.To,
		ToInclusive: r.ToInclusive,
	}
}

func recordFromKey(table, key string) string {
	return "type::record('" + table + "', " + key + ")"
}
