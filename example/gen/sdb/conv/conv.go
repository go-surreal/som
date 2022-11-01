package conv

import (
  "github.com/google/uuid"
	"strings"
	"time"
)

func prepareID(node string, id any) string {
	return strings.TrimPrefix(id.(string), node+":")
}

func parseTime(val any) time.Time {
	res, err := time.Parse(time.RFC3339, val.(string))
	if err != nil {
		return time.Time{}
	}
	return res
}

func parseUUID(val any) uuid.UUID {
	res, err := uuid.Parse(val.(string))
	if err != nil {
		return uuid.UUID{}
	}
	return res
}
	
// func extract[T any](val any, to func(map[string]any) T) T {
//	var t T
//	if val == nil {
//		return t
//	}
//	return to(val.(map[string]any))
// }
