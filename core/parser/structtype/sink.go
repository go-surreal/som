package structtype

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type SinkHandler struct{}

func (h *SinkHandler) Match(t gotype.Type, ctx *parser.TypeContext) bool {
	return IsSink(t, ctx.OutPkg)
}

func (h *SinkHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	sink, err := ParseSink(t, ctx.OutPkg)
	if err != nil {
		return err
	}
	ctx.Output.Sinks = append(ctx.Output.Sinks, sink)
	return nil
}

func (h *SinkHandler) Validate(ctx *parser.TypeContext) error {
	names := make([]string, len(ctx.Output.Sinks))
	for i, s := range ctx.Output.Sinks {
		names[i] = s.Name
		if err := validateFields("sink "+s.Name, s.Fields, ctx.Output); err != nil {
			return err
		}
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate sink name %q", dup)
	}
	return nil
}

func IsSink(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if !f.IsAnonymous() {
			continue
		}

		if f.Name() == "Sink" && f.Elem().Name() == "Sink" &&
			f.Elem().PkgPath() == outPkg {
			return true
		}
	}

	return false
}

func ParseSink(v gotype.Type, outPkg string) (*parser.Sink, error) {
	sink := &parser.Sink{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "Sink" {
				// The sink marker (som.Sink) contributes no column.
				continue
			}

			return nil, fmt.Errorf("sink %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("sink %s: field ID not allowed, a sink record is discarded and has no addressable id", v.Name())
		}

		field, err := parser.ParseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		sink.Fields = append(sink.Fields, field)
	}

	return sink, nil
}
