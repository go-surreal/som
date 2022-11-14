package conv

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

func parseDatabaseID(node string, id string) string {
	return strings.TrimPrefix(id, node+":")
}

func buildDatabaseID(node string, id string) string {
	return node + ":" + id
}

func parseTime(val any) time.Time {
	res, err := time.Parse(time.RFC3339, val.(string))
	if err != nil {
		return time.Time{}
	}
	return res
}

func parseUUID(val string) uuid.UUID {
	res, err := uuid.Parse(val)
	if err != nil {
		// TODO: add logging!
		return uuid.UUID{}
	}
	return res
}
	
func mapRecords[I, O any](in []I, fn func(*I) O) []O {
	var out []O
	for _, i := range in {
		out = append(out, fn(&i))
	}
	return out
}
	
func convertEnum[I, O ~string](in []I) []O {
	var out []O
	for _, i := range in {
		out = append(out, O(i))
	}
	return out
}

// func extract[T any](val any, to func(map[string]any) T) T {
//	var t T
//	if val == nil {
//		return t
//	}
//	return to(val.(map[string]any))
// }
