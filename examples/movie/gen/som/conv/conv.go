// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package conv

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func parseDatabaseID(node string, id string) string {
	id = strings.TrimPrefix(id, node+":")
	id = strings.TrimPrefix(id,"⟨")
	id = strings.TrimSuffix(id, "⟩")
	id, _ = strconv.Unquote("\"" + id + "\"")
	return id
}

func buildDatabaseID(node string, id string) string {
	return node + ":" + id
}

func parseTime(val any) time.Time {
	res, err := time.Parse(time.RFC3339, val.(string))
	if err != nil {
		// TODO: add logging!
		return time.Time{}
	}
	return res
}

func durationPtr(val *time.Duration) *string {
	if val == nil {
		return nil
	}
	str := val.String()
	return &str
}

func parseDuration(val string) time.Duration {
	res, err := time.ParseDuration(val)
	if err != nil {
		// TODO: add logging!
		return 0
	}
	return res
}

func uuidPtr(val *uuid.UUID) *string {
	if val == nil {
		return nil
	}
	str := val.String()
	return &str
}

func parseUUID(val string) uuid.UUID {
	res, err := uuid.Parse(val)
	if err != nil {
		// TODO: add logging!
		return uuid.UUID{}
	}
	return res
}
	
func mapEnum[I, O ~string](in I) O {
 	return O(in)
}

func mapSlice[I, O any](in []I, fn func(I) O) []O {
	if in == nil {
		return nil
	}

	out := make([]O, len(in))
	for _, i := range in {
		out = append(out, fn(i))
	}
	return out
}
	
func mapSlicePtr[I, O any](in *[]I, fn func(I) O) *[]O {
	if in == nil {
		return nil
	}

	out := make([]O, len(*in))
	for _, i := range *in {
		out = append(out, fn(i))
	}
	return &out
}

func mapPtrSlice[I, O any](in []*I, fn func(I) O) []*O {
	if in == nil {
		return nil
	}

 	ptrFn := ptrFunc(fn)

	out := make([]*O, len(in))
 	for _, i := range in {
 		out = append(out, ptrFn(i))
 	}

 	return out
}

func mapPtrSlicePtr[I, O any](in *[]*I, fn func(I) O) *[]*O {
	if in == nil {
		return nil
	}

	ptrFn := ptrFunc(fn)

	out := make([]*O, len(*in))
	for _, i := range *in {
		out = append(out, ptrFn(i))
	}

	return &out
}
	
func ptrFunc[I, O any](fn func(I) O) func(*I) *O {
 	return func(in *I) *O {
 		if in == nil {
 			return nil
 		}
 		out := fn(*in)
 		return &out
 	}
}
