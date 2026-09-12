package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// Enum is a string type declared in the model package. It is also its own
// TypeHandler.
type Enum struct {
	Name string
}

func (e *Enum) Match(t gotype.Type, ctx *TypeContext) bool {
	return IsEnum(t, ctx.OutPkg)
}

func (e *Enum) Handle(t gotype.Type, ctx *TypeContext) error {
	ctx.Output.Enums = append(ctx.Output.Enums, &Enum{
		Name: t.Name(),
	})
	return nil
}

func (e *Enum) Validate(ctx *TypeContext) error {
	names := make([]string, len(ctx.Output.Enums))
	for i, enum := range ctx.Output.Enums {
		names[i] = enum.Name
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate enum name %q", dup)
	}
	return nil
}

func IsEnum(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.String {
		return false
	}

	return t.String() != "string" && t.PkgPath() == outPkg
}
