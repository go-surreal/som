package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Password struct {
	*baseField

	source *parser.FieldPassword
}

func (f *Password) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, "Password").Types(jen.Qual(f.SourcePkg, string(f.source.Algorithm)))
}

func (f *Password) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(ctx.TargetPkg, "Password").Types(jen.Qual(ctx.TargetPkg, string(f.source.Algorithm)))
}

func (f *Password) TypeDatabase() string {
	return f.optionWrap("string")
}

func (f *Password) cryptoGenerateFunc() string {
	switch f.source.Algorithm {
	case parser.PasswordBcrypt:
		return "crypto::bcrypt::generate"
	case parser.PasswordArgon2:
		return "crypto::argon2::generate"
	case parser.PasswordPbkdf2:
		return "crypto::pbkdf2::generate"
	case parser.PasswordScrypt:
		return "crypto::scrypt::generate"
	default:
		return "crypto::bcrypt::generate"
	}
}

func (f *Password) SchemaStatements(table, prefix string) []string {
	// Only hash if value is present AND different from $before (prevents double-hashing on updates)
	valueClause := fmt.Sprintf(
		`IF $value != NONE AND $value != NULL AND $value != "" AND $value != $before THEN %s($value) ELSE $value END`,
		f.cryptoGenerateFunc(),
	)

	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s VALUE %s PERMISSIONS FOR SELECT NONE;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(), valueClause,
		),
	}
}

func (f *Password) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,

		fieldDef: f.fieldDef,
	}
}

func (f *Password) filterDefine(ctx Context) jen.Code {
	filter := "Password"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Password) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewPassword"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(
			jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
			jen.Qual(ctx.pkgLib(), string(f.source.Algorithm)).Values(),
		)
}

func (f *Password) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Password) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}

	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *Password) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
