package parser

// FieldVersion is the revision counter contributed by the som.OptimisticLock
// embed. It is never matched against a source field, the model type creates it.
type FieldVersion struct {
	fieldBase
}

func NewFieldVersion(name string) *FieldVersion {
	return &FieldVersion{fieldBase: newBase(name)}
}
