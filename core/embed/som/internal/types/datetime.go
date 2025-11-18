//go:build embed

package types

import (
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
)

const (
	tagDatetime         = 12
	expectedArrayLength = 2
)

// DateTime wraps time.Time and provides CBOR marshaling for SurrealDB.
type DateTime struct {
	time.Time
}

func (dt *DateTime) MarshalCBOR() ([]byte, error) {
	if dt == nil {
		data, err := cbor.Marshal(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal nil: %w", err)
		}

		return data, nil
	}

	content, err := cbor.Marshal([]int64{dt.Unix(), int64(dt.Nanosecond())})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal datetime slice: %w", err)
	}

	data, err := cbor.Marshal(cbor.RawTag{
		Number:  tagDatetime,
		Content: content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw tag: %w", err)
	}

	return data, nil
}

func (dt *DateTime) UnmarshalCBOR(data []byte) error {
	var val []int64

	if err := cbor.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("failed to unmarshal datetime: %w", err)
	}

	if len(val) == 0 {
		// Empty array means zero/unset time
		dt.Time = time.Time{}
		return nil
	}

	if len(val) > expectedArrayLength {
		return fmt.Errorf("invalid datetime array length: expected 0-2 elements, got %d", len(val))
	}

	secs := val[0]
	nano := int64(0)

	if len(val) > 1 {
		nano = val[1]
	}

	dt.Time = time.Unix(secs, nano)

	return nil
}
