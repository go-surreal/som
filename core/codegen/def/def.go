package def

import "github.com/dave/jennifer/jen"

const (
	PkgQuery  = "query"
	PkgFilter = "filter"
	PkgSort   = "by"
	PkgFetch  = "with"
	PkgConv   = "conv"
	PkgRelate = "relate"
	PkgField  = "field"
	PkgRepo    = "repo"
	PkgSomWire = "somwire"

	PkgInternal    = "internal"
	PkgLib         = "internal/lib"
	PkgTypes       = "internal/types"
	PkgCBORHelpers = "internal/cbor"
	PkgDistinct    = "internal/distinct"

	IndexPrefix = "__som__"
)

const (
	//PkgSom  = "github.com/go-surreal/som"
	PkgSurrealDB = "github.com/surrealdb/surrealdb.go"
	PkgModels    = "github.com/surrealdb/surrealdb.go/pkg/models"
	PkgCBOR      = "github.com/fxamacker/cbor/v2"

	PkgURL        = "net/url"
	PkgUUIDGoogle = "github.com/google/uuid"
	PkgUUIDGofrs  = "github.com/gofrs/uuid"
)

var (
	TypeModel = jen.Id("M")
)
