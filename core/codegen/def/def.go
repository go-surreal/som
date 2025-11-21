package def

import "github.com/dave/jennifer/jen"

const (
	PkgQuery  = "query"
	PkgFilter = "where"
	PkgSort   = "by"
	PkgFetch  = "with"
	PkgConv   = "conv"
	PkgRelate = "relate"
	PkgRepo   = "repo"

	PkgLib   = "internal/lib"
	PkgTypes = "internal/types"
	PkgCBORHelpers = "internal/cbor"
)

const (
	//PkgSom  = "github.com/go-surreal/som"
	PkgSurrealDB = "github.com/surrealdb/surrealdb.go"
	PkgModels    = "github.com/surrealdb/surrealdb.go/pkg/models"
	PkgCBOR      = "github.com/fxamacker/cbor/v2"

	PkgURL  = "net/url"
	PkgUUID = "github.com/google/uuid"
)

var (
	TypeModel = jen.Id("M")
)
