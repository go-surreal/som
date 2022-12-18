package field

import (
	"github.com/marcbinz/som/core/parser"
)

func Convert(conf *BuildConfig, field parser.Field, getElement ElemGetter) (Field, bool) {
	base := &baseField{BuildConfig: conf, source: field}

	switch f := field.(type) {

	case *parser.FieldID:
		return &ID{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldString:
		return &String{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldNumeric:
		return &Numeric{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldBool:
		return &Bool{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldTime:
		return &Time{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldUUID:
		return &UUID{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldNode:
		return &Node{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldEnum:
		return &Enum{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldStruct:
		return &Struct{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldSlice:
		return &Slice{
			baseField:  base,
			source:     f,
			getElement: getElement,
		}, true
	}

	return nil, false
}
