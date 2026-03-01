package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/paulmach/orb"
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
}
