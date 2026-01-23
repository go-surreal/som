//go:build embed

package lib

import (
	"strconv"
	"strings"
	"time"
)

type ChangesQuery struct {
	Table string
	Since any // time.Time or uint64 (versionstamp)
	Limit int
}

func (q ChangesQuery) Build() *Result {
	var out strings.Builder

	out.WriteString("SHOW CHANGES FOR TABLE ")
	out.WriteString(q.Table)
	out.WriteString(" SINCE ")

	switch v := q.Since.(type) {
	case time.Time:
		// SurrealDB requires datetime literals in the format: d"2023-09-07T01:23:52Z"
		// Parameters are not supported for SHOW CHANGES SINCE.
		out.WriteString(`d"`)
		out.WriteString(v.UTC().Format(time.RFC3339Nano))
		out.WriteString(`"`)
	case uint64:
		out.WriteString(strconv.FormatUint(v, 10))
	default:
		out.WriteString("0")
	}

	if q.Limit > 0 {
		out.WriteString(" LIMIT ")
		out.WriteString(strconv.Itoa(q.Limit))
	}

	return &Result{
		Statement: out.String(),
		Variables: nil,
	}
}
