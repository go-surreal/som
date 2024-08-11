//go:build embed

package conv

import (
	"encoding/json"
	"github.com/go-surreal/sdbc"
	"github.com/go-surreal/som/core/embed/internal/types"
	"github.com/google/uuid"
	"net/url"
	"strconv"
	"strings"
	"time"

	"{{.GenerateOutPath}}/internal/types"
)

func parseDatabaseID(node string, id string) string {
	id = strings.TrimPrefix(id, node+":")
	id = strings.TrimPrefix(id, "⟨")
	id = strings.TrimSuffix(id, "⟩")
	id, _ = strconv.Unquote("\"" + id + "\"")
	return id
}

func buildDatabaseID(node string, id string) string {
	return node + ":" + id
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

func mapSliceFn[I, O any](fn func(I) O) func([]I) []O {
	return func(in []I) []O {
		if in == nil {
			return nil
		}

		out := make([]O, len(in))

		for _, i := range in {
			out = append(out, fn(i))
		}

		return out
	}
}

func mapSliceFnPtr[I, O any](fn func(I) O) func(*[]I) *[]O {
	return func(in *[]I) *[]O {
		if in == nil {
			return nil
		}

		out := make([]O, len(*in))

		for _, i := range *in {
			out = append(out, fn(i))
		}

		return &out
	}
}

func noOp[I any](in I) I {
	return in
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

func noPtrFunc[I, O any](fn func(*I) *O) func(I) O {
	return func(in I) O {
		out := fn(&in)
		if out == nil {
			var o O
			return o
		}
		return *out
	}
}

//
// -- TIME
//

func fromTime(val time.Time) sdbc.DateTime {
	return sdbc.DateTime{val}
}

func toTime(val sdbc.DateTime) time.Time {
	return val.Time
}

func fromTimePtr(val *time.Time) *sdbc.DateTime {
	if val == nil {
		return nil
	}

	return &sdbc.DateTime{*val}
}

func toTimePtr(val *sdbc.DateTime) *time.Time {
	if val == nil {
		return nil
	}

	return &val.Time
}

//
// -- NUMBER
//

type unsignedNumber[T uint | uint64 | uintptr] struct {
	val *T
}

func (n *unsignedNumber[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(*n.val), 10) + "dec")
}

func (n *unsignedNumber[T]) UnmarshalJSON(data []byte) error {
	if n == nil {
		return nil
	}

	var raw string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if raw == "" {
		return nil
	}

	res, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return err
	}

	val := T(res)
	n.val = &val

	return nil
}

//
// -- URL
//

func fromURL(val url.URL) string {
	return val.String()
}

func fromURLPtr(val *url.URL) *string {
	if val == nil {
		return nil
	}

	str := val.String()
	return &str
}

func toURL(val string) url.URL {
	res, err := url.Parse(val)
	if err != nil {
		// TODO: add logging!
		return url.URL{}
	}

	return *res
}

func toURLPtr(val *string) *url.URL {
	if val == nil {
		return nil
	}

	res, err := url.Parse(*val)
	if err != nil {
		// TODO: add logging!
		return nil
	}

	return res
}

//
// -- DURATION
//

func fromDuration(val time.Duration) sdbc.Duration {
	return sdbc.Duration{val}
}

func fromDurationPtr(val *time.Duration) *sdbc.Duration {
	if val == nil {
		return nil
	}

	return &sdbc.Duration{*val}
}

func toDuration(val sdbc.Duration) time.Duration {
	return val.Duration
}

func toDurationPtr(val *sdbc.Duration) *time.Duration {
	if val == nil {
		return nil
	}

	return &val.Duration
}

//
// -- UUID
//

func fromUUID(val uuid.UUID) types.UUID {
	return types.UUID(val)
}

func fromUUIDPtr(val *uuid.UUID) *types.UUID {
	if val == nil {
		return nil
	}

	u := types.UUID(*val)
	return &u
}

func toUUID(val types.UUID) uuid.UUID {
	return uuid.UUID(val)
}

func toUUIDPtr(val *types.UUID) *uuid.UUID {
	if val == nil {
		return nil
	}

	u := uuid.UUID(*val)
	return &u
}
