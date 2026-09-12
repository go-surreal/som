package parser

type IDType string

const (
	IDTypeULID   IDType = "ULID"
	IDTypeUUID   IDType = "UUID"
	IDTypeRand   IDType = "Rand"
	IDTypeString IDType = "String"
	IDTypeArray  IDType = "Array"
	IDTypeObject IDType = "Object"
)

// FieldID is the record id contributed by the som.Node or som.Edge embed. It
// is never matched against a source field, the model type creates it.
type FieldID struct {
	fieldBase
	Type IDType
}

func NewFieldID(name string, idType IDType) *FieldID {
	return &FieldID{fieldBase: newBase(name), Type: idType}
}

func (f *FieldID) IsInternal() bool {
	return true
}
