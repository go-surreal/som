package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

type TypeContext struct {
	OutPkg   string
	PkgScope gotype.Type
	Output   *Output
}

type TypeHandler interface {
	Priority() int
	Match(t gotype.Type, ctx *TypeContext) bool
	Handle(t gotype.Type, ctx *TypeContext) error
	Validate(ctx *TypeContext) error
}

var defaultTypeRegistry = &typeRegistry{}

type typeRegistry struct {
	registry[TypeHandler]
}

func RegisterTypeHandler(h TypeHandler) {
	defaultTypeRegistry.register(h)
}

func (r *typeRegistry) handle(t gotype.Type, ctx *TypeContext) (bool, error) {
	for _, h := range r.all(func(h TypeHandler) int { return h.Priority() }) {
		if h.Match(t, ctx) {
			return true, h.Handle(t, ctx)
		}
	}
	return false, nil
}

func (r *typeRegistry) validate(ctx *TypeContext) error {
	for _, h := range r.all(func(h TypeHandler) int { return h.Priority() }) {
		if err := h.Validate(ctx); err != nil {
			return fmt.Errorf("validation: %w", err)
		}
	}
	return nil
}
