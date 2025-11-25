//go:build embed

package lib

import (
	"strings"
)

// CursorFilter creates a compound filter for cursor-based pagination.
// It generates conditions like:
// ((a > $cursor_a) OR (a = $cursor_a AND b < $cursor_b) OR (a = $cursor_a AND b = $cursor_b AND id > $cursor_id))
func CursorFilter[M any](cursor CursorData, sorts []*SortBuilder, backward bool) Filter[M] {
	return filter[M](func(ctx *context, _ M) string {
		if len(sorts) == 0 {
			return ""
		}

		// Build compound OR conditions for multi-field cursor comparison
		var orParts []string

		for i := range sorts {
			var andParts []string

			// Add equality conditions for all previous fields
			for j := 0; j < i; j++ {
				field := sorts[j].Field
				val, ok := cursor.SortValues[field]
				if !ok && field == "id" {
					val = cursor.ID
				}
				if val == nil {
					continue
				}
				varName := ctx.asVar(val)
				andParts = append(andParts, field+" = "+varName)
			}

			// Add comparison condition for current field
			field := sorts[i].Field
			val, ok := cursor.SortValues[field]
			if !ok && field == "id" {
				val = cursor.ID
			}
			if val == nil {
				continue
			}

			varName := ctx.asVar(val)
			op := cursorComparisonOp(sorts[i].Order, backward)
			andParts = append(andParts, field+" "+string(op)+" "+varName)

			if len(andParts) == 1 {
				orParts = append(orParts, andParts[0])
			} else {
				orParts = append(orParts, "("+strings.Join(andParts, " AND ")+")")
			}
		}

		if len(orParts) == 0 {
			return ""
		}

		if len(orParts) == 1 {
			return orParts[0]
		}

		return "(" + strings.Join(orParts, " OR ") + ")"
	})
}

// cursorComparisonOp returns the comparison operator for cursor pagination.
// For forward pagination: ASC uses >, DESC uses <
// For backward pagination: ASC uses <, DESC uses >
func cursorComparisonOp(order SortOrder, backward bool) Operator {
	if backward {
		if order == SortAsc {
			return OpLessThan
		}
		return OpGreaterThan
	}
	// Forward
	if order == SortAsc {
		return OpGreaterThan
	}
	return OpLessThan
}
