//go:build embed

package conv

import (
	"encoding/json"
	"github.com/fxamacker/cbor/v2"
	"github.com/go-surreal/sdbc"
	"github.com/google/uuid"
	"net/url"
	"strconv"
	"strings"
	"time"
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

// func uintToString[T uint | uint64 | uintptr](val T) string {
// 	return strconv.FormatUint(uint64(val), 10)
// }
//
// func stringToUint[T uint | uint64 | uintptr](val string) T {
// 	res, err := strconv.ParseUint(val, 10, 64)
// 	if err != nil {
// 		return 0
// 	}
// 	return T(res)
// }

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

func urlPtr(val *url.URL) *string {
	if val == nil {
		return nil
	}
	str := val.String()
	return &str
}

func parseURL(val string) url.URL {
	res, err := url.Parse(val)
	if err != nil {
		// TODO: add logging!
		return url.URL{}
	}
	return *res
}

//
// -- UUID
//

type UUID uuid.UUID

func (u *UUID) MarshalCBOR() ([]byte, error) {
	if u == nil {
		return nil, nil
	}

	raw, err := cbor.Marshal(uuid.UUID(*u))
	if err != nil {
		return nil, err
	}

	return cbor.Marshal(cbor.RawTag{
		Number:  sdbc.CBORTagUUID,
		Content: raw,
	})
}
