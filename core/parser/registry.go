package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

type registry[H any] struct {
	handlers []H
}

// Type registry

type TypeContext struct {
	OutPkg   string
	PkgScope gotype.Type
	Output   *Output
}

type TypeHandler interface {
	Match(t gotype.Type, ctx *TypeContext) bool
	Handle(t gotype.Type, ctx *TypeContext) error
	Validate(ctx *TypeContext) error
}

type typeRegistry struct {
	registry[TypeHandler]
}

func newTypeRegistry(handlers []TypeHandler) *typeRegistry {
	return &typeRegistry{registry[TypeHandler]{handlers: handlers}}
}

func (r *typeRegistry) handle(t gotype.Type, ctx *TypeContext) (bool, error) {
	for _, h := range r.handlers {
		if h.Match(t, ctx) {
			return true, h.Handle(t, ctx)
		}
	}
	return false, nil
}

func (r *typeRegistry) validate(ctx *TypeContext) error {
	for _, h := range r.handlers {
		if err := h.Validate(ctx); err != nil {
			return fmt.Errorf("validation: %w", err)
		}
	}
	return nil
}

// Field registry

type FieldContext struct {
	OutPkg       string
	RecurseParse func(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error)
}

type FieldHandler interface {
	Match(elem gotype.Type, ctx *FieldContext) bool
	Parse(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error)
}

type fieldRegistry struct {
	registry[FieldHandler]
}

func newFieldRegistry(handlers []FieldHandler) *fieldRegistry {
	return &fieldRegistry{registry[FieldHandler]{handlers: handlers}}
}

func (r *fieldRegistry) parse(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error) {
	for _, h := range r.handlers {
		if h.Match(elem, ctx) {
			return h.Parse(t, elem, ctx)
		}
	}
	return nil, fmt.Errorf("field %s has unsupported type %s", t.Name(), elem.Kind())
}
