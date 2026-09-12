package parser

import (
	"fmt"
	"go/ast"

	"github.com/wzshiming/gotype"
)

// Struct is a plain nested struct stored as an embedded object.
//
// It is also its own TypeHandler, and since it matches any struct it must be
// registered last, after every more specific model type.
type Struct struct {
	Name   string
	Fields []Field
}

func (s *Struct) Match(t gotype.Type, _ *TypeContext) bool {
	return t.Kind() == gotype.Struct
}

func (s *Struct) Handle(t gotype.Type, ctx *TypeContext) error {
	str, err := ParseStruct(t, ctx)
	if err != nil {
		return err
	}
	ctx.Output.Structs = append(ctx.Output.Structs, str)
	return nil
}

func (s *Struct) Validate(ctx *TypeContext) error {
	names := make([]string, len(ctx.Output.Structs))
	for i, str := range ctx.Output.Structs {
		names[i] = str.Name
		if err := validateFields("struct "+str.Name, str.Fields, ctx.Output); err != nil {
			return err
		}
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate struct name %q", dup)
	}
	return nil
}

func ParseStruct(v gotype.Type, ctx *TypeContext) (*Struct, error) {
	str := &Struct{Name: v.Name()}

	nf := v.NumField()

	for i := range nf {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		field, err := ctx.ParseField(f)
		if err != nil {
			return nil, fmt.Errorf("struct %s: %w", v.Name(), err)
		}

		str.Fields = append(str.Fields, field)
	}

	return str, nil
}
