package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/wzshiming/gotype"
)

// Sink is a write-only ingestion table backed by a DEFINE TABLE ... DROP
// statement. It has a struct shape (its columns) like a Node but no ID,
// no features and only create operations; rows are discarded immediately
// after write, so they cannot be read back or linked to.
//
// It is also its own TypeHandler.
type Sink struct {
	Name   string
	Fields []Field
}

func (s *Sink) Match(t gotype.Type, ctx *TypeContext) bool {
	return IsSink(t, ctx.OutPkg)
}

func (s *Sink) Handle(t gotype.Type, ctx *TypeContext) error {
	sink, err := ParseSink(t, ctx)
	if err != nil {
		return err
	}
	ctx.Output.Sinks = append(ctx.Output.Sinks, sink)
	return nil
}

func (s *Sink) Validate(ctx *TypeContext) error {
	names := make([]string, len(ctx.Output.Sinks))
	for i, sink := range ctx.Output.Sinks {
		names[i] = sink.Name
		if err := validateFields("sink "+sink.Name, sink.Fields, ctx.Output); err != nil {
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

	for i := range nf {
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

func ParseSink(v gotype.Type, ctx *TypeContext) (*Sink, error) {
	outPkg := ctx.OutPkg

	sink := &Sink{Name: v.Name()}

	nf := v.NumField()

	for i := range nf {
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

		field, err := ctx.ParseField(f)
		if err != nil {
			return nil, err
		}

		sink.Fields = append(sink.Fields, field)
	}

	return sink, nil
}
