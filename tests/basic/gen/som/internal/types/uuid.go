// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package types

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/go-surreal/sdbc"
	"github.com/google/uuid"
)

type UUID uuid.UUID

func (u *UUID) MarshalCBOR() ([]byte, error) {
	if u == nil {
		return cbor.Marshal(nil)
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