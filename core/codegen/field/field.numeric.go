package field

import (
	"fmt"
	"math"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
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
	//case parser.NumberUint:
	//	return jen.Add(f.ptr()).Uint()
	case parser.NumberUint8:
		return jen.Add(f.ptr()).Uint8()
	case parser.NumberUint16:
		return jen.Add(f.ptr()).Uint16()
	case parser.NumberUint32:
		return jen.Add(f.ptr()).Uint32()
	//case parser.NumberUint64:
	//	return jen.Add(f.ptr()).Uint64()
	//case parser.NumberUintptr:
	//	return jen.Add(f.ptr()).Uintptr()
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

func (f *Numeric) typeConv(_ Context) jen.Code {
	switch f.source.Type {

	//case parser.NumberUint64, parser.NumberUint, parser.NumberUintptr:
	//	{
	//		var typ jen.Code
	//		//// Qual(def.PkgJson, "Number")
	//		switch f.source.Type {
	//		case parser.NumberUint:
	//			typ = jen.Uint()
	//		case parser.NumberUint64:
	//			typ = jen.Uint64()
	//		case parser.NumberUintptr:
	//			typ = jen.Uintptr()
	//		}
	//
	//		return jen.Id("unsignedNumber").Types(typ)
	//	}

	default:
		return f.typeGo()
	}
}

func (f *Numeric) TypeDatabase() string {
	switch f.source.Type {
	case parser.NumberInt8, parser.NumberInt16, parser.NumberInt32, parser.NumberRune,
		parser.NumberInt64, parser.NumberInt,
		parser.NumberUint8, parser.NumberUint16, parser.NumberUint32:
		return f.optionWrap("int")
	case parser.NumberFloat32, parser.NumberFloat64:
		return f.optionWrap("float")
	default:
		panic(fmt.Sprintf("unmapped numeric type: %d", f.source.Type))
	}
}

func (f *Numeric) SchemaStatements(table, prefix string) []string {
	extend := f.typeDatabaseExtend()
	return []string{f.schemaStatement(table, prefix, f.TypeDatabase(), extend)}
}

func (f *Numeric) typeDatabaseExtend() string {
	nilCheck := ""
	if f.source.Pointer() {
		nilCheck = "$value == NONE OR $value == NULL OR "
	}

	switch f.source.Type {
	case parser.NumberInt8:
		return fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, math.MinInt8, math.MaxInt8)
	case parser.NumberInt16:
		return fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, math.MinInt16, math.MaxInt16)
	case parser.NumberInt32, parser.NumberRune:
		return fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, math.MinInt32, math.MaxInt32)
	case parser.NumberInt64, parser.NumberInt:
		return fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, math.MinInt64, math.MaxInt64)
	case parser.NumberUint8:
		return fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, 0, math.MaxUint8)
	case parser.NumberUint16:
		return fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, 0, math.MaxUint16)
	case parser.NumberUint32:
		return fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, 0, math.MaxUint32)
	case parser.NumberFloat32, parser.NumberFloat64:
		return ""
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

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
		fieldDef:      f.fieldDef,
	}
}

func (f *Numeric) filterDefine(ctx Context) jen.Code {
	filter := "Numeric"

	switch f.source.Type {

	case parser.NumberInt, parser.NumberInt8, parser.NumberInt16, parser.NumberInt32, parser.NumberInt64,
		parser.NumberUint8, parser.NumberUint16, parser.NumberUint32, parser.NumberRune:
		{
			filter = "Int"
		}

	case parser.NumberFloat32, parser.NumberFloat64:
		{
			filter = "Float"
		}
	}

	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel, f.typeGo())
}

func (f *Numeric) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewNumeric"

	switch f.source.Type {

	case parser.NumberInt, parser.NumberInt8, parser.NumberInt16, parser.NumberInt32, parser.NumberInt64,
		parser.NumberUint8, parser.NumberUint16, parser.NumberUint32, parser.NumberRune:
		{
			filter = "NewInt"
		}

	case parser.NumberFloat32, parser.NumberFloat64:
		{
			filter = "NewFloat"
		}
	}

	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel, f.typeGo()),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Numeric) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Numeric) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Numeric) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Numeric) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *Numeric) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
