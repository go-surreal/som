package def

import "github.com/dave/jennifer/jen"

const (
	PkgQuery  = "query"
	PkgFilter = "where"
	PkgSort   = "by"
	PkgFetch  = "with"
	PkgConv   = "conv"
	PkgRelate = "relate"

	PkgLib   = "internal/lib"
	PkgTypes = "internal/types"
)

const (
	//PkgSom  = "github.com/go-surreal/som"
	PkgSDBC = "github.com/go-surreal/sdbc"
	PkgCBOR = "github.com/fxamacker/cbor/v2"

	PkgURL  = "net/url"
	PkgUUID = "github.com/google/uuid"
)

var (
	TypeModel = jen.Id("M")
)
