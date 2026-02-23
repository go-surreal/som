package parser

import (
	"go/ast"

	"github.com/wzshiming/gotype"
)

type passwordFieldHandler struct{}



func (h *passwordFieldHandler) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.String && isPassword(elem, ctx.OutPkg)
}

func (h *passwordFieldHandler) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldPassword{&fieldAtomic{name: t.Name()}, parsePasswordAlgorithm(elem)}, nil
}

func isPassword(t gotype.Type, outPkg string) bool {
	if t.PkgPath() != outPkg {
		return false
	}
	return t.Name() == "Password"
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
