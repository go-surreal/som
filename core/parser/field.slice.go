package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// FieldSlice is a slice of any other supported field type.
type FieldSlice struct {
	fieldBase
	Field Field
}

func (f *FieldSlice) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Slice
}

func (f *FieldSlice) Parse(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error) {
	inner, err := ctx.ParseElem(elem, elem.Elem())
	if err != nil {
		return nil, err
	}

	return &FieldSlice{fieldBase: newBase(t.Name()), Field: inner}, nil
}

func (f *FieldSlice) Validate() error {
	if err := f.Field.Validate(); err != nil {
		return err
	}

	if f.search != nil {
		if _, ok := f.Field.(*FieldString); !ok {
			return fmt.Errorf("field %s: fulltext index only supports string slice types, got slice of %T", f.name, f.Field)
		}
	}

	if union, ok := f.Field.(*FieldUnion); ok {
		return fmt.Errorf("field %s: a slice of the union %s is not supported yet", f.name, union.Union)
	}

	return nil
}

func (f *FieldSlice) Verify(out *Output) error {
	if err := f.fieldBase.Verify(out); err != nil {
		return err
	}

	return f.Field.Verify(out)
}

func (f *FieldSlice) contributeFeatures(features *UsedFeatures) {
	if c, ok := f.Field.(featureContributor); ok {
		c.contributeFeatures(features)
	}
}
