package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// FieldEnum is a string type declared in the model package, used as an enum.
type FieldEnum struct {
	fieldBase
	Typ string
}

func (f *FieldEnum) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.String && IsEnum(elem, ctx.OutPkg)
}

func (f *FieldEnum) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldEnum{fieldBase: newBase(t.Name()), Typ: elem.Name()}, nil
}

func (f *FieldEnum) Verify(out *Output) error {
	if err := f.fieldBase.Verify(out); err != nil {
		return err
	}

	if !enumExists(f.Typ, out) {
		return fmt.Errorf("field %q references unknown enum %q", f.Name(), f.Typ)
	}

	return nil
}
