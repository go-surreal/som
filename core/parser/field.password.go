package parser

import (
	"go/ast"

	"github.com/wzshiming/gotype"
)

type PasswordAlgorithm string

const (
	PasswordBcrypt PasswordAlgorithm = "Bcrypt"
	PasswordArgon2 PasswordAlgorithm = "Argon2"
	PasswordPbkdf2 PasswordAlgorithm = "Pbkdf2"
	PasswordScrypt PasswordAlgorithm = "Scrypt"
)

// FieldPassword is a som.Password, parameterised with its hash algorithm.
type FieldPassword struct {
	fieldBase
	Algorithm PasswordAlgorithm
}

func (f *FieldPassword) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.String && elem.PkgPath() == ctx.OutPkg && elem.Name() == "Password"
}

func (f *FieldPassword) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldPassword{fieldBase: newBase(t.Name()), Algorithm: parsePasswordAlgorithm(elem)}, nil
}

func parsePasswordAlgorithm(t gotype.Type) PasswordAlgorithm {
	origin := t.Origin()
	if origin == nil {
		return PasswordBcrypt
	}

	if indexExpr, ok := origin.(*ast.IndexExpr); ok {
		if selExpr, ok := indexExpr.Index.(*ast.SelectorExpr); ok {
			switch selExpr.Sel.Name {
			case "Bcrypt":
				return PasswordBcrypt
			case "Argon2":
				return PasswordArgon2
			case "Pbkdf2":
				return PasswordPbkdf2
			case "Scrypt":
				return PasswordScrypt
			}
		}
	}

	return PasswordBcrypt
}
