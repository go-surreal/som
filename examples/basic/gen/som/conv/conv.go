// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package conv

import (
	"encoding/json"
	"fmt"
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
		return time.Time{}
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

//
// -- UUID
//

type UUID struct {
	*uuid.UUID
}

func (u *UUID) MarshalJSON() ([]byte, error) {
	if u == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(u.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("cannot unmarshal uuid: %w", err)
	}

	u.UUID = &uid

	return nil
}
