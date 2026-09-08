package field

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

type NodeTable struct {
	Name       string
	Fields     []Field
	Changefeed string
	Source     *parser.Node // Reference to source parser.Node

	// TODO: include source package path + method(s)
}

func (t *NodeTable) FileName() string {
	return "node." + strcase.ToSnake(t.Name) + ".go"
}

func (t *NodeTable) GetFields() []Field {
	return t.Fields
}

func (t *NodeTable) NameGo() string {
	return t.Name
}

func (t *NodeTable) NameGoLower() string {
	return strcase.ToLowerCamel(t.Name)
}

func (t *NodeTable) NameDatabase() string {
	return strcase.ToSnake(t.Name) // TODO
}

func (t *NodeTable) HasChangefeed() bool {
	return t.Changefeed != ""
}

func (t *NodeTable) HasComplexID() bool {
	return t.Source != nil && t.Source.ComplexID != nil
}

func (t *NodeTable) HasStringID() bool {
	if t.Source == nil {
		return false
	}
	switch t.Source.IDType {
	case parser.IDTypeULID, parser.IDTypeUUID, parser.IDTypeRand, parser.IDTypeString:
		return true
	}
	return false
}

// HasAutoID reports whether record IDs for this table are generated
// server-side when none is provided. It is false for som.String and for
// complex (array/object) IDs, which must always be supplied explicitly.
func (t *NodeTable) HasAutoID() bool {
	if t.Source == nil {
		return false
	}
	switch t.Source.IDType {
	case parser.IDTypeULID, parser.IDTypeUUID, parser.IDTypeRand:
		return true
	}
	return false
}
