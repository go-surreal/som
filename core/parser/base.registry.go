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

	fields *fieldRegistry
}

// ParseField parses a struct field of a model type, including its som tag.
func (c *TypeContext) ParseField(t gotype.Type) (Field, error) {
	tagInfo, err := parseSomTag(t.Tag().Get("som"))
	if err != nil {
		return nil, fmt.Errorf("field %s: %w", t.Name(), err)
	}

	field, err := c.parseFieldType(t)
	if err != nil {
		return nil, err
	}

	if tagInfo != nil {
		if tagInfo.DBName != "" {
			field.setDBName(tagInfo.DBName)
		}
		if len(tagInfo.Indexes) > 0 {
			field.setIndexes(tagInfo.Indexes)
		}
		if tagInfo.Search != nil {
			field.setSearch(tagInfo.Search)
		}
	}

	if err := field.Validate(); err != nil {
		return nil, err
	}

	return field, nil
}

// ParseUntaggedField parses a field that carries no som tag of its own, e.g. a
// field of a complex ID key struct.
func (c *TypeContext) ParseUntaggedField(t gotype.Type) (Field, error) {
	return c.parseFieldType(t)
}

func (c *TypeContext) parseFieldType(t gotype.Type) (Field, error) {
	ctx := &FieldContext{OutPkg: c.OutPkg}
	ctx.ParseElem = func(t gotype.Type, elem gotype.Type) (Field, error) {
		return c.fields.parse(t, elem, ctx)
	}

	return c.fields.parse(t, t.Elem(), ctx)
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
	OutPkg    string
	ParseElem func(t gotype.Type, elem gotype.Type) (Field, error)
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
	isPtr := elem.Kind() == gotype.Ptr
	if isPtr {
		elem = elem.Elem()
	}

	for _, h := range r.handlers {
		if h.Match(elem, ctx) {
			field, err := h.Parse(t, elem, ctx)
			if err != nil {
				return nil, err
			}
			if isPtr {
				field.setPointer(true)
			}
			return field, nil
		}
	}
	return nil, fmt.Errorf("field %s has unsupported type %s", t.Name(), elem.Kind())
}
