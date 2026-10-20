package parser

import (
	"fmt"
	"go/ast"
	"path"
	"strings"

	"github.com/wzshiming/gotype"
)

// Edge is a relation table between two nodes. It is also its own TypeHandler.
type Edge struct {
	Name           string
	In             Field
	Out            Field
	Fields         []Field
	Timestamps     bool
	OptimisticLock bool
	Changefeed     string
	SoftDelete     bool
}

func (e *Edge) Match(t gotype.Type, ctx *TypeContext) bool {
	return IsEdge(t, ctx.OutPkg)
}

func (e *Edge) Handle(t gotype.Type, ctx *TypeContext) error {
	edge, err := ParseEdge(t, ctx)
	if err != nil {
		return err
	}
	ctx.Output.Edges = append(ctx.Output.Edges, edge)
	return nil
}

func (e *Edge) Validate(ctx *TypeContext) error {
	names := make([]string, len(ctx.Output.Edges))
	for i, edge := range ctx.Output.Edges {
		names[i] = edge.Name
		if edge.In == nil {
			return fmt.Errorf("edge %s: missing 'in' field (tag som:\"in\")", edge.Name)
		}
		if edge.Out == nil {
			return fmt.Errorf("edge %s: missing 'out' field (tag som:\"out\")", edge.Name)
		}
		if err := edge.In.Verify(ctx.Output); err != nil {
			return fmt.Errorf("edge %s: %w", edge.Name, err)
		}
		if err := edge.Out.Verify(ctx.Output); err != nil {
			return fmt.Errorf("edge %s: %w", edge.Name, err)
		}
		if err := validateFields("edge "+edge.Name, edge.Fields, ctx.Output); err != nil {
			return err
		}
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate edge name %q", dup)
	}
	return nil
}

func IsEdge(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := range nf {
		f := t.Field(i)

		if !f.IsAnonymous() {
			continue
		}

		if f.Name() == "Edge" && f.Elem().Name() == "Edge" &&
			f.Elem().PkgPath() == outPkg {
			return true
		}
	}

	return false
}

func ParseEdge(v gotype.Type, ctx *TypeContext) (*Edge, error) {
	outPkg := ctx.OutPkg
	internalPkg := path.Join(outPkg, "internal")

	edge := &Edge{Name: v.Name()}

	var features FeatureSet

	nf := v.NumField()

	for i := range nf {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "Edge" {
				edge.Changefeed = ParseChangefeedTag(f.Tag().Get("som"))
				edge.Fields = append(edge.Fields,
					NewFieldID("ID", IDTypeULID),
				)
				continue
			}

			matched, err := ParseFeature(f, internalPkg, &features, &edge.Fields)
			if err != nil {
				return nil, fmt.Errorf("edge %s: %w", v.Name(), err)
			}
			if matched {
				if features.Expiry {
					return nil, fmt.Errorf("edge %s: som.Expiry is not supported on edges", v.Name())
				}
				continue
			}

			return nil, fmt.Errorf("edge %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("edge %s: field ID not allowed, already provided by som.Edge", v.Name())
		}

		field, err := ctx.ParseField(f)
		if err != nil {
			return nil, fmt.Errorf("edge %s: %w", v.Name(), err)
		}

		if f.Tag().Get("som") == "in" {
			edge.In = field
			continue
		}

		if f.Tag().Get("som") == "out" {
			edge.Out = field
			continue
		}

		edge.Fields = append(edge.Fields, field)
	}

	ApplyFeatures(features, &edge.Timestamps, &edge.OptimisticLock, &edge.SoftDelete, &edge.Fields)

	return edge, nil
}
