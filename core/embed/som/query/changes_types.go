//go:build embed

package query

import "github.com/fxamacker/cbor/v2"

// ChangeEntry represents a batch of changes at a specific versionstamp.
type ChangeEntry[M any] struct {
	Versionstamp uint64
	Creates      []M
	Updates      []M
	Deletes      []M
}

// rawChangeEntry is the raw response structure from SHOW CHANGES.
type rawChangeEntry struct {
	Versionstamp uint64      `cbor:"versionstamp"`
	Changes      []rawChange `cbor:"changes"`
}

// rawChange represents a single change operation.
// Only one of the fields will be set.
type rawChange struct {
	DefineTable *cbor.RawMessage `cbor:"define_table,omitempty"`
	Update      cbor.RawMessage  `cbor:"update,omitempty"`
	Create      cbor.RawMessage  `cbor:"create,omitempty"`
	Delete      cbor.RawMessage  `cbor:"delete,omitempty"`
}
