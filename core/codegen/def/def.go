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

	PkgLib         = "internal/lib"
	PkgTypes       = "internal/types"
	PkgCBORHelpers = "internal/cbor"

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

	PkgGeoOrb            = "github.com/paulmach/orb"
	PkgGeoSimplefeatures = "github.com/peterstace/simplefeatures/geom"
)

var (
	TypeModel = jen.Id("M")
)
