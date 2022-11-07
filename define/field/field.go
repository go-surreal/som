package field

type dataType int32

const (
	typeOther dataType = iota
	typeString
	typeInt
	typeInt32
	typeInt64
	typeFloat32
	typeFloat64
	typeBool
	typeUUID
	typeTime
	typeLink
	typeObject
)

type Field struct {
	name      string
	dataType  dataType
	isPointer bool
	isSlice   bool
	isUnique  bool
	asLink    TableDef
	asObject  ObjectDef
}

func (f *Field) Pointer() *Field {
	f.isPointer = true
	return f
}

func (f *Field) Slice() *Field {
	f.isSlice = true
	return f
}

func (f *Field) Unique() *Field {
	f.isUnique = true
	return f
}

// func (f *Field) Index(idx any) *Field {
// 	return f
// }

// func (f *Field) Value(val *value.Value) *Field {
// 	return f
// }
//
// func (f *Field) Assert(assert any) *Field {
// 	return f
// }

func String(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeString,
	}
}

func Int(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeInt,
	}
}

func Int32(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeInt32,
	}
}

func Int64(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeInt64,
	}
}

func Float32(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeFloat32,
	}
}

func Float64(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeFloat64,
	}
}

func Bool(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeBool,
	}
}

func UUID(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeUUID,
	}
}

func Time(name string) *Field {
	return &Field{
		name:     name,
		dataType: typeTime,
	}
}

func Object(name string, object ObjectDef) *Field {
	return &Field{
		name:     name,
		dataType: typeObject,
		asObject: object,
	}
}

func Link(name string, table TableDef) *Field {
	return &Field{
		name:     name,
		dataType: typeLink,
		asLink:   table,
	}
}

type TableDef interface{}

type ObjectDef interface{}
