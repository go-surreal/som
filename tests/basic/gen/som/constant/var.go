package constant

import "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"

func String[M any](val string) *lib.String[M] {
	return lib.NewString[M](lib.NewVarKey[M](val))
}
