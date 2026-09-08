package def

import "github.com/dave/jennifer/jen"

const (
	PkgQuery  = "query"
	PkgFilter = "filter"
	PkgSort   = "by"
	PkgFetch  = "with"
	PkgConv   = "conv"
	PkgRelate = "relate"
	PkgIndex  = "index"
	PkgField  = "field"
	PkgRepo    = "repo"
	PkgSomWire = "somwire"

	PkgInternal    = "internal"
	PkgLib         = "internal/lib"
	PkgTypes       = "internal/types"
	PkgCBORHelpers = "internal/cbor"
	PkgDistinct    = "internal/distinct"

	IndexPrefix = "__som__"

	// MetaTable stores the hash of each expensive schema object (index, view),
	// so that re-applying an unchanged schema does not rebuild it.
	MetaTable = "__som__schema"
)

const (
	//PkgSom  = "github.com/go-surreal/som"
	PkgSurrealDB = "github.com/surrealdb/surrealdb.go"
	PkgModels    = "github.com/surrealdb/surrealdb.go/pkg/models"

	PkgURL        = "net/url"
	PkgUUIDGoogle = "github.com/google/uuid"
	PkgUUIDGofrs  = "github.com/gofrs/uuid"
	PkgUUIDStd    = "uuid"

	PkgGeoOrb            = "github.com/paulmach/orb"
	PkgGeoSimplefeatures = "github.com/peterstace/simplefeatures/geom"
	PkgGeoGoGeom         = "github.com/twpayne/go-geom"
)

var (
	TypeModel = jen.Id("M")
)
