package fieldtype

import (
	"go/ast"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type PasswordHandler struct{}

func (h *PasswordHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.String && isPassword(elem, ctx.OutPkg)
}

func (h *PasswordHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldPassword(t.Name(), parsePasswordAlgorithm(elem)), nil
}

func isPassword(t gotype.Type, outPkg string) bool {
	if t.PkgPath() != outPkg {
		return false
	}
	return t.Name() == "Password"
}

func parsePasswordAlgorithm(t gotype.Type) parser.PasswordAlgorithm {
	origin := t.Origin()
	if origin == nil {
		return parser.PasswordBcrypt
	}

	if indexExpr, ok := origin.(*ast.IndexExpr); ok {
		if selExpr, ok := indexExpr.Index.(*ast.SelectorExpr); ok {
			switch selExpr.Sel.Name {
			case "Bcrypt":
				return parser.PasswordBcrypt
			case "Argon2":
				return parser.PasswordArgon2
			case "Pbkdf2":
				return parser.PasswordPbkdf2
			case "Scrypt":
				return parser.PasswordScrypt
			}
		}
	}

	return parser.PasswordBcrypt
}
