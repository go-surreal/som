package parser

import (
	"fmt"
	"go/ast"
	"path"
	"strings"

	"github.com/wzshiming/gotype"
)

type edgeHandler struct{}



func (h *edgeHandler) Match(t gotype.Type, ctx *TypeContext) bool {
	return isEdge(t, ctx.OutPkg)
}

func (h *edgeHandler) Handle(t gotype.Type, ctx *TypeContext) error {
	edge, err := parseEdge(t, ctx.OutPkg)
	if err != nil {
		return err
	}
	ctx.Output.Edges = append(ctx.Output.Edges, edge)
	return nil
}

func (h *edgeHandler) Validate(_ *TypeContext) error { return nil }

func isEdge(t gotype.Type, outPkg string) bool {
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

func parseEdge(v gotype.Type, outPkg string) (*Edge, error) {
	internalPkg := path.Join(outPkg, "internal")

	edge := &Edge{Name: v.Name()}

	var features featureSet

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "Edge" {
				edge.Fields = append(edge.Fields,
					&FieldID{&fieldAtomic{name: "ID"}, IDTypeULID},
				)
				continue
			}

			if parseFeature(f, internalPkg, &features, &edge.Fields) {
				continue
			}

			return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.Edge", v.Name())
		}

		field, err := parseField(f, outPkg)
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

	applyFeatures(features, &edge.Timestamps, &edge.OptimisticLock, &edge.SoftDelete, &edge.Fields)

	return edge, nil
}
