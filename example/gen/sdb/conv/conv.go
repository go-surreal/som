package conv

import (
	uuid "github.com/google/uuid"
	"time"
)

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
