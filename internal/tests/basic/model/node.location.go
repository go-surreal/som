package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/paulmach/orb"
	sfgeom "github.com/peterstace/simplefeatures/geom"
	"github.com/twpayne/go-geom"
)

type Location struct {
	som.Node[som.ULID]
	som.Timestamps

	Name string

	Point           orb.Point
	PointPtr        *orb.Point
	LineString      orb.LineString
	LineStringPtr   *orb.LineString
	Polygon         orb.Polygon
	PolygonPtr      *orb.Polygon
	MultiPoint      orb.MultiPoint
	MultiPointPtr   *orb.MultiPoint
	MultiLineString orb.MultiLineString
	MultiPolygon    orb.MultiPolygon
	Collection      orb.Collection

	GGPoint           geom.Point
	GGPointPtr        *geom.Point
	GGLineString      geom.LineString
	GGLineStringPtr   *geom.LineString
	GGPolygon         geom.Polygon
	GGPolygonPtr      *geom.Polygon
	GGMultiPoint      geom.MultiPoint
	GGMultiPointPtr   *geom.MultiPoint
	GGMultiLineString geom.MultiLineString
	GGMultiPolygon    geom.MultiPolygon

	SFPoint           sfgeom.Point
	SFPointPtr        *sfgeom.Point
	SFLineString      sfgeom.LineString
	SFLineStringPtr   *sfgeom.LineString
	SFPolygon         sfgeom.Polygon
	SFPolygonPtr      *sfgeom.Polygon
	SFMultiPoint      sfgeom.MultiPoint
	SFMultiPointPtr   *sfgeom.MultiPoint
	SFMultiLineString sfgeom.MultiLineString
	SFMultiPolygon    sfgeom.MultiPolygon
}
