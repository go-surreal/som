package field

import (
	"bytes"
	"github.com/dave/jennifer/jen"
)

func isCodeEqual(a, b jen.Code) bool {
	if (a == nil || b == nil) && a != b {
		return false
	}

	aBytes := bytes.NewBuffer(nil)
	bBytes := bytes.NewBuffer(nil)

	if err := jen.Add(a).Render(aBytes); err != nil {
		return false
	}

	if err := jen.Add(b).Render(bBytes); err != nil {
		return false
	}

	return bytes.Equal(aBytes.Bytes(), bBytes.Bytes())
}
