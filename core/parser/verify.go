package parser

import (
	"fmt"
)

// Verify performs the post-parse checks of a field, i.e. everything that can
// only be decided once all types of the model package are known. It is the
// counterpart to Validate, which runs while the field itself is parsed and
// therefore only sees the field in isolation.
//
// The returned error is not prefixed with the owning type, the caller adds
// that context.
func (f *fieldAtomic) Verify(out *Output) error {
	if f.search != nil && !searchExists(f.search.ConfigName, out) {
		return fmt.Errorf("field %q references unknown search config %q", f.name, f.search.ConfigName)
	}

	return nil
}

// IsInternal reports whether the field is managed by som itself and may
// therefore occupy a reserved database name.
func (f *fieldAtomic) IsInternal() bool {
	return false
}

func (f *FieldID) IsInternal() bool {
	return true
}

func (f *FieldTime) IsInternal() bool {
	return f.IsCreatedAt || f.IsUpdatedAt || f.IsDeletedAt || f.IsExpiresAt
}

func (f *FieldNode) Verify(out *Output) error {
	if err := f.fieldAtomic.Verify(out); err != nil {
		return err
	}

	if !nodeExists(f.Node, out) {
		return fmt.Errorf("field %q references unknown node %q", f.Name(), f.Node)
	}

	return nil
}

func (f *FieldEdge) Verify(out *Output) error {
	if err := f.fieldAtomic.Verify(out); err != nil {
		return err
	}

	if !edgeExists(f.Edge, out) {
		return fmt.Errorf("field %q references unknown edge %q", f.Name(), f.Edge)
	}

	return nil
}

func (f *FieldEnum) Verify(out *Output) error {
	if err := f.fieldAtomic.Verify(out); err != nil {
		return err
	}

	if !enumExists(f.Typ, out) {
		return fmt.Errorf("field %q references unknown enum %q", f.Name(), f.Typ)
	}

	return nil
}

func (f *FieldStruct) Verify(out *Output) error {
	if err := f.fieldAtomic.Verify(out); err != nil {
		return err
	}

	if !structExists(f.Struct, out) {
		return fmt.Errorf("field %q references unknown struct %q", f.Name(), f.Struct)
	}

	return nil
}

func (f *FieldSlice) Verify(out *Output) error {
	if err := f.fieldAtomic.Verify(out); err != nil {
		return err
	}

	return f.Field.Verify(out)
}

func (f *FieldComplexID) Verify(out *Output) error {
	if err := f.fieldAtomic.Verify(out); err != nil {
		return err
	}

	for _, sf := range f.Fields {
		if err := sf.Field.Verify(out); err != nil {
			return err
		}
	}

	return nil
}

func (f *FieldComplexID) IsInternal() bool {
	return true
}

func nodeExists(name string, out *Output) bool {
	for _, n := range out.Nodes {
		if n.Name == name {
			return true
		}
	}
	return false
}

func edgeExists(name string, out *Output) bool {
	for _, e := range out.Edges {
		if e.Name == name {
			return true
		}
	}
	return false
}

func enumExists(name string, out *Output) bool {
	for _, e := range out.Enums {
		if e.Name == name {
			return true
		}
	}
	return false
}

func structExists(name string, out *Output) bool {
	for _, s := range out.Structs {
		if s.Name == name {
			return true
		}
	}
	return false
}

func searchExists(name string, out *Output) bool {
	if out.Define == nil {
		return false
	}
	for _, s := range out.Define.Searches {
		if s.Name == name {
			return true
		}
	}
	return false
}
