package constant

import "{{.GenerateOutPath}}/internal/lib"

func String[M any](val string) *lib.String[M] {
	return lib.NewString[M](lib.NewVarKey[M](val))
}
