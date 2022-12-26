package som

import (
	"time"
)

type Record = Node // TODO: should we use this?

type Node struct {
	// include query info into each node resulting from a query?:
	// status string
	// time   string
	// extract via som.Info(someNode) -> som.Info ?
}

// Edge describes an edge between two Node elements.
// It may have its own fields.
type Edge struct{}

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Enum describes a database type with a fixed set of allowed values.
type Enum string

// Password describes a special string field.
// Regarding the generated database query operations, it can only be matched, but never read.
// In a query result, the Password field will always be empty.
type Password string

// Meta describes a model that is not related to any Node or Edge.
// Instead, it is used to hold metadata that was queried from a Node or Edge.
//
// Applying this struct to a type within your model package will ensure
// that this type is never considered for the generated layer.
type Meta struct{}

// Info holds information about a single database operation.
// It is used as target to hold said information when building
// an operation using the WithInfo() method. The generated
// builder for each model provides this capability.
//
// Example:
// Take a model named "Model" for which the som code is generated.
// Accessing the database operation happens as usual for example via
// client.Model().Create() or client.Model().Query(). Extracting the
// Info out of those operations is as simple as:
//
// var info *som.Info
// client.Model().WithInfo(info).Create()
// client.Model().WithInfo(info).Query()
//
// Please note: When using the same base for multiple operations, the Info
// struct will only ever hold the information of the last operation.
type Info struct {
	Time    time.Time
	Status  string
	Message string
}

type Entity interface {
	entity()
}

//
// -- GEO
//

type Point struct {
	Longitude float64
	Latitude  float64
}

// {
//    type: "Point",
//    coordinates: [-0.118092, 51.509865],
// }

type Line struct {
	Points []Point
}

// {
//	type: "LineString",
//	coordinates: [
//		[10.0, 11.2], [10.5, 11.9]
//	]
// }

type Polygon struct {
	Points []MultiPoint
}

// {
//	type: "Polygon",
//	coordinates: [[
//		[-0.38314819, 51.37692386], [0.1785278, 51.37692386],
//		[0.1785278, 51.61460570], [-0.38314819, 51.61460570],
//		[-0.38314819, 51.37692386]
//	]]
// }

type MultiPoint struct {
	Points []Point
}

// {
//	type: "MultiPoint",
//	coordinates: [
//		[10.0, 11.2],
//		[10.5, 11.9]
//	],
// }

type MultiLine struct {
	Lines []Line
}

// {
//	type: "MultiLinestring",
//	coordinates: [
//		[ [10.0, 11.2], [10.5, 11.9] ],
//		[ [11.0, 12.2], [11.5, 12.9], [12.0, 13.0] ]
//	]
// }

type MultiPolygon struct {
	Polygons []Polygon
}

// {
//	type: "MultiPolygon",
//	coordinates: [
//		[
//			[ [10.0, 11.2], [10.5, 11.9], [10.8, 12.0], [10.0, 11.2] ]
//		],
//		[
//			[ [9.0, 11.2], [10.5, 11.9], [10.3, 13.0], [9.0, 11.2] ]
//		]
//	]
// }

type Collection struct {
	Geometries []any
}

// {
//	type: "GeometryCollection",
//	geometries: [
//		{
//			type: "MultiPoint",
//			coordinates: [
//				[10.0, 11.2],
//				[10.5, 11.9]
//			],
//		},
//		{
//			type: "Polygon",
//			coordinates: [[
//				[-0.38314819, 51.37692386], [0.1785278, 51.37692386],
//				[0.1785278, 51.61460570], [-0.38314819, 51.61460570],
//				[-0.38314819, 51.37692386]
//			]]
//		},
//		{
//			type: "MultiPolygon",
//			coordinates: [
//				[
//					[ [10.0, 11.2], [10.5, 11.9], [10.8, 12.0], [10.0, 11.2] ]
//				],
//				[
//					[ [9.0, 11.2], [10.5, 11.9], [10.3, 13.0], [9.0, 11.2] ]
//				]
//			]
//		}
//	]
// }
