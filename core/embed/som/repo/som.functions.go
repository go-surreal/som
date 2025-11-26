//go:build embed

package repo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

const statusOK = "OK"

type RawQuery[T any] struct {
	Status string
	Time   string
	Result T
	Detail string
}

func Unmarshal[M any](respond any) (model M, err error) {
	var bytes []byte

	if arrResp, isArr := respond.([]any); len(arrResp) > 0 {
		if dataMap, ok := arrResp[0].(map[string]any); ok && isArr {
			if _, ok := dataMap["status"]; ok {
				if bytes, err = json.Marshal(respond); err == nil {
					var raw []RawQuery[M]
					if err = json.Unmarshal(bytes, &raw); err == nil {
						if raw[0].Status != statusOK {
							err = fmt.Errorf("%s: %s", raw[0].Status, raw[0].Detail)
						}
						model = raw[0].Result
					}
				}
				return model, err
			}
		}
	}

	if bytes, err = json.Marshal(respond); err == nil {
		err = json.Unmarshal(bytes, &model)
	}

	return model, err
}

//
// -- RECORD ID
//

const recordSeparator = ":"

var ErrUnmarshalNotSupported = errors.New("unmarshal not supported")

type RecordID interface {
	recordID()
}

type newRecordID struct {
	table       string
	constructor string
}

func (id *newRecordID) recordID() {}

func (id *newRecordID) String() string {
	return id.table + recordSeparator + id.constructor
}

func (id *newRecordID) MarshalCBOR() ([]byte, error) {
	content, err := cbor.Marshal(id.table + recordSeparator + id.constructor)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recordID: %w", err)
	}

	data, err := cbor.Marshal(cbor.RawTag{
		Number:  models.TagRecordID,
		Content: content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw tag: %w", err)
	}

	return data, nil
}

func (id *newRecordID) UnmarshalCBOR(_ []byte) error {
	return ErrUnmarshalNotSupported
}

func newID(table string) RecordID {
	return &newRecordID{
		table:       table,
		constructor: "rand()",
	}
}

func newULID(table string) RecordID {
	return &newRecordID{
		table:       table,
		constructor: "ulid()",
	}
}

func newUUID(table string) RecordID {
	return &newRecordID{
		table:       table,
		constructor: "uuid()", // TODO: schema type for ID field and cbor tag
	}
}
