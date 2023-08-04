package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/parser"
)

type Numeric struct {
	*baseField

	source *parser.FieldNumeric
}

func (f *Numeric) typeGo() jen.Code {
	switch f.source.Type {
	case parser.NumberInt:
		return jen.Add(f.ptr()).Int()
	case parser.NumberInt8:
		return jen.Add(f.ptr()).Int8()
	case parser.NumberInt16:
		return jen.Add(f.ptr()).Int16()
	case parser.NumberInt32:
		return jen.Add(f.ptr()).Int32()
	case parser.NumberInt64:
		return jen.Add(f.ptr()).Int64()
	case parser.NumberUint:
		return jen.Add(f.ptr()).Uint()
	case parser.NumberUint8:
		return jen.Add(f.ptr()).Uint8()
	case parser.NumberUint16:
		return jen.Add(f.ptr()).Uint16()
	case parser.NumberUint32:
		return jen.Add(f.ptr()).Uint32()
	case parser.NumberUint64:
		return jen.Add(f.ptr()).Uint64()
	case parser.NumberUintptr:
		return jen.Add(f.ptr()).Uintptr()
	case parser.NumberFloat32:
		return jen.Add(f.ptr()).Float32()
	case parser.NumberFloat64:
		return jen.Add(f.ptr()).Float64()
	case parser.NumberRune:
		return jen.Add(f.ptr()).Rune()
	default:
		panic(fmt.Sprintf("unmapped numeric type: %d", f.source.Type))
	}
}

func (f *Numeric) typeConv() jen.Code {
	return f.typeGo()
}

func (f *Numeric) TypeDatabase() string {
	nilCheck := ""
	if f.source.Pointer() {
		nilCheck = "$value == NONE OR $value == NULL OR "
	}

	switch f.source.Type {
	case parser.NumberInt8:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= -128 AND $value <= 127"
	case parser.NumberInt16:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= -32768 AND $value <= 32767"
	case parser.NumberInt32, parser.NumberRune:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= -2147483648 AND $value <= 2147483647"
	case parser.NumberInt64, parser.NumberInt:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= -9223372036854775808 AND $value <= 9223372036854775807"
	case parser.NumberUint8:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= 0 AND $value <= 255"
	case parser.NumberUint16:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= 0 AND $value <= 65535"
	case parser.NumberUint32:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= 0 AND $value <= 4294967295"
	case parser.NumberUint64, parser.NumberUint, parser.NumberUintptr:
		return f.optionWrap("int") + " ASSERT " + nilCheck + "$value >= 0 AND $value <= 18446744073709551615"
	case parser.NumberFloat32:
		return f.optionWrap("float") + " ASSERT " + nilCheck + "$value >= -3.402823466E+38 AND $value <= 3.402823466E+38"
	case parser.NumberFloat64:
		return f.optionWrap("float") + " ASSERT " + nilCheck + "$value >= -1.7976931348623157E+308 AND $value <= 1.7976931348623157E+308"
	default:
		panic(fmt.Sprintf("unmapped numeric type: %d", f.source.Type))
	}
}

func (f *Numeric) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Numeric) filterDefine(ctx Context) jen.Code {
	filter := "Numeric"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(f.typeGo(), jen.Id("T"))
}

func (f *Numeric) filterInit(ctx Context) jen.Code {
	filter := "NewNumeric"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(f.typeGo(), jen.Id("T")).
		Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Numeric) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(jen.Id("T"))
}

func (f *Numeric) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Numeric) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Numeric) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Numeric) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
