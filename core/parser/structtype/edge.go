package structtype

import (
	"fmt"
	"go/ast"
	"path"
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type EdgeHandler struct{}

func (h *EdgeHandler) Match(t gotype.Type, ctx *parser.TypeContext) bool {
	return IsEdge(t, ctx.OutPkg)
}

func (h *EdgeHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	edge, err := ParseEdge(t, ctx.OutPkg)
	if err != nil {
		return err
	}
	ctx.Output.Edges = append(ctx.Output.Edges, edge)
	return nil
}

func (h *EdgeHandler) Validate(_ *parser.TypeContext) error { return nil }

func IsEdge(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
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

func ParseEdge(v gotype.Type, outPkg string) (*parser.Edge, error) {
	internalPkg := path.Join(outPkg, "internal")

	edge := &parser.Edge{Name: v.Name()}

	var features parser.FeatureSet

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "Edge" {
				edge.Fields = append(edge.Fields,
					parser.NewFieldID("ID", parser.IDTypeULID),
				)
				continue
			}

			if parser.ParseFeature(f, internalPkg, &features, &edge.Fields) {
				continue
			}

			return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.Edge", v.Name())
		}

		field, err := parser.ParseField(f, outPkg)
		if err != nil {
			return nil, err
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

	parser.ApplyFeatures(features, &edge.Timestamps, &edge.OptimisticLock, &edge.SoftDelete, &edge.Fields)

	return edge, nil
}
