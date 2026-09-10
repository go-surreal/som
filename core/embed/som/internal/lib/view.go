//go:build embed

package lib

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// RenderLiteralFilter renders a filter as a static SurrealQL condition with
// all values inlined as literals (no $-parameters). It is used to build the
// WHERE clause of a DEFINE TABLE ... AS SELECT view definition, which cannot
// carry query parameters.
func RenderLiteralFilter[M any](filters ...Filter[M]) string {
	if len(filters) == 0 {
		return ""
	}

	ctx := &context{vars: map[string]any{}, literal: true}

	var m M
	out := All[M](filters).build(ctx, m)

	// All wraps a single condition in parentheses; keep the raw form for a
	// WHERE clause and let the caller decide on grouping.
	out = strings.TrimPrefix(out, "(")
	out = strings.TrimSuffix(out, ")")

	return out
}

// literalValue renders a Go value as a SurrealQL literal.
func literalValue(val any) string {
	switch v := val.(type) {
	case nil:
		return "NONE"
	case bool:
		return strconv.FormatBool(v)
	case string:
		return "'" + strings.ReplaceAll(v, "'", "\\'") + "'"
	case int:
		return strconv.Itoa(v)
	case int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		// Any slice/array (e.g. the typed []string / []int passed by IN /
		// NotIn filters) renders as a SurrealQL array literal.
		//
		// TODO: this does not yet render converted values faithfully (e.g.
		// internal/types.DateTime -> d"...", Duration, RecordID). Such values
		// currently fall through to the %v fallback and produce invalid DDL;
		// a proper literal renderer for view WHERE clauses is still needed.
		if rv := reflect.ValueOf(val); rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			parts := make([]string, rv.Len())
			for i := range parts {
				parts[i] = literalValue(rv.Index(i).Interface())
			}
			return "[" + strings.Join(parts, ", ") + "]"
		}
		return fmt.Sprintf("%v", v)
	}
}
