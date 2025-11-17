//go:build embed

package types

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/surrealdb/surrealdb.go/pkg/models"
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
		Number:  models.TagSpecBinaryUUID,
		Content: raw,
	})
}

func (u *UUID) UnmarshalCBOR(data []byte) error {
	var tag cbor.RawTag
	if err := cbor.Unmarshal(data, &tag); err != nil {
		return err
	}

	if tag.Number != models.TagSpecBinaryUUID {
		return fmt.Errorf("unexpected tag number for UUID: got %d, want %d", tag.Number, models.TagSpecBinaryUUID)
	}

	// tag.Content is cbor.RawMessage which is []byte
	// We need to unmarshal it to get the actual UUID bytes
	var uuidBytes []byte
	if err := cbor.Unmarshal(tag.Content, &uuidBytes); err != nil {
		return fmt.Errorf("failed to unmarshal UUID content: %w", err)
	}

	if len(uuidBytes) != 16 {
		return fmt.Errorf("UUID must be exactly 16 bytes, got %d", len(uuidBytes))
	}

	parsed, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return fmt.Errorf("failed to parse UUID: %w", err)
	}

	*u = UUID(parsed)
	return nil
}
