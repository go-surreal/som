package parser

import (
	"github.com/wzshiming/gotype"
)

// validateMethodName is the method a type may declare to have its values
// checked before they are written to the database.
const validateMethodName = "Validate"

// hasValidateMethod reports whether the given type provides a
// "Validate() error" method, either declared on the type itself or promoted
// from an embedded one. A pointer receiver counts as well, since every value
// som validates is addressable.
func hasValidateMethod(t gotype.Type) bool {
	if t == nil {
		return false
	}

	if t.Kind() == gotype.Declaration {
		t = t.Declaration()
	}

	method, ok := t.MethodByName(validateMethodName)
	if !ok {
		return false
	}

	if method.Kind() == gotype.Declaration {
		method = method.Declaration()
	}

	if method.Kind() != gotype.Func {
		return false
	}

	if method.NumIn() != 0 || method.NumOut() != 1 {
		return false
	}

	return isErrorType(method.Out(0))
}

// isErrorType reports whether the given type is the builtin error interface.
// A result is reported as a declaration, named or not, so it is unwrapped
// before the kind is checked.
func isErrorType(t gotype.Type) bool {
	if t == nil {
		return false
	}

	if t.Kind() == gotype.Declaration {
		t = t.Declaration()
	}

	return t.Kind() == gotype.Error
}
