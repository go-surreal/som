package basic

import (
	"context"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/paulmach/orb"
	sfgeom "github.com/peterstace/simplefeatures/geom"
	"github.com/twpayne/go-geom"
	"gotest.tools/v3/assert"
)

// minimalLocation returns a Location with valid minimal data for all non-pointer geo fields.
// SurrealDB requires valid geometry data even for "empty" geometries.
func minimalLocation(name string) *model.Location {
	return &model.Location{
		Name:            name,
		Point:           orb.Point{0, 0},
		LineString:      orb.LineString{{0, 0}, {1, 1}},
		Polygon:         orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
		MultiPoint:      orb.MultiPoint{{0, 0}},
		MultiLineString: orb.MultiLineString{{{0, 0}, {1, 1}}},
		MultiPolygon:    orb.MultiPolygon{{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}}},
		Collection:      orb.Collection{orb.Point{0, 0}},

		GGPoint:           *geom.NewPointFlat(geom.XY, []float64{0, 0}),
		GGLineString:      *geom.NewLineStringFlat(geom.XY, []float64{0, 0, 1, 1}),
		GGPolygon:         *geom.NewPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, []int{5}),
		GGMultiPoint:      *geom.NewMultiPointFlat(geom.XY, []float64{0, 0}),
		GGMultiLineString: *geom.NewMultiLineStringFlat(geom.XY, []float64{0, 0, 1, 1}, []int{2}),
		GGMultiPolygon:    *geom.NewMultiPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, [][]int{{5}}),

		SFPoint:           sfNewPoint(0, 0),
		SFLineString:      sfNewLineString([]float64{0, 0, 1, 1}),
		SFPolygon:         sfNewPolygon([][]float64{{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}}),
		SFMultiPoint:      sfgeom.NewMultiPoint([]sfgeom.Point{sfNewPoint(0, 0)}),
		SFMultiLineString: sfgeom.NewMultiLineString([]sfgeom.LineString{sfNewLineString([]float64{0, 0, 1, 1})}),
		SFMultiPolygon:    sfgeom.NewMultiPolygon([]sfgeom.Polygon{sfNewPolygon([][]float64{{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}})}),
	}
}

func sfNewPoint(x, y float64) sfgeom.Point {
	return sfgeom.NewPoint(sfgeom.Coordinates{XY: sfgeom.XY{X: x, Y: y}})
}

func sfNewLineString(flat []float64) sfgeom.LineString {
	return sfgeom.NewLineString(sfgeom.NewSequence(flat, sfgeom.DimXY))
}

func sfNewPolygon(rings [][]float64) sfgeom.Polygon {
	lsRings := make([]sfgeom.LineString, len(rings))
	for i, ring := range rings {
		lsRings[i] = sfNewLineString(ring)
	}
	return sfgeom.NewPolygon(lsRings)
}

func TestGeoPoint(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test Point")
	loc.Point = orb.Point{1.0, 2.0}
	ptr := orb.Point{2.0, 3.0}
	loc.PointPtr = &ptr

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Timestamps{},
			geom.Point{}, geom.LineString{}, geom.Polygon{},
			geom.MultiPoint{}, geom.MultiLineString{}, geom.MultiPolygon{},
			sfgeom.Point{}, sfgeom.LineString{}, sfgeom.Polygon{},
			sfgeom.MultiPoint{}, sfgeom.MultiLineString{}, sfgeom.MultiPolygon{},
		),
	)

	assert.Equal(t, loc.Point[0], out.Point[0])
	assert.Equal(t, loc.Point[1], out.Point[1])
	assert.Equal(t, (*loc.PointPtr)[0], (*out.PointPtr)[0])
	assert.Equal(t, (*loc.PointPtr)[1], (*out.PointPtr)[1])
}

func TestGeoLineString(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test LineString")
	loc.LineString = orb.LineString{{1.0, 2.0}, {3.0, 4.0}, {5.0, 6.0}}
	ptr := orb.LineString{{3.0, 4.0}, {5.0, 6.0}}
	loc.LineStringPtr = &ptr

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Timestamps{},
			geom.Point{}, geom.LineString{}, geom.Polygon{},
			geom.MultiPoint{}, geom.MultiLineString{}, geom.MultiPolygon{},
			sfgeom.Point{}, sfgeom.LineString{}, sfgeom.Polygon{},
			sfgeom.MultiPoint{}, sfgeom.MultiLineString{}, sfgeom.MultiPolygon{},
		),
	)

	assert.Equal(t, len(loc.LineString), len(out.LineString))
	for i := range loc.LineString {
		assert.Equal(t, loc.LineString[i][0], out.LineString[i][0])
		assert.Equal(t, loc.LineString[i][1], out.LineString[i][1])
	}
}

func TestGeoPolygon(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test Polygon")

	// Simple square polygon (exterior ring only)
	exterior := orb.Ring{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	loc.Polygon = orb.Polygon{exterior}

	// Polygon with a hole
	hole := orb.Ring{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}
	polygonWithHole := orb.Polygon{exterior, hole}
	loc.PolygonPtr = &polygonWithHole

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Timestamps{},
			geom.Point{}, geom.LineString{}, geom.Polygon{},
			geom.MultiPoint{}, geom.MultiLineString{}, geom.MultiPolygon{},
			sfgeom.Point{}, sfgeom.LineString{}, sfgeom.Polygon{},
			sfgeom.MultiPoint{}, sfgeom.MultiLineString{}, sfgeom.MultiPolygon{},
		),
	)

	// Verify polygon structure
	assert.Equal(t, len(loc.Polygon), len(out.Polygon))
	assert.Equal(t, len((*loc.PolygonPtr)), len((*out.PolygonPtr)))
}

func TestGeoMultiTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test Multi Types")

	loc.MultiPoint = orb.MultiPoint{{1, 2}, {3, 4}, {5, 6}}
	loc.MultiLineString = orb.MultiLineString{
		{{0, 0}, {1, 1}},
		{{2, 2}, {3, 3}, {4, 4}},
	}
	loc.MultiPolygon = orb.MultiPolygon{
		{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
		{{{2, 2}, {3, 2}, {3, 3}, {2, 3}, {2, 2}}},
	}

	multiPointPtr := orb.MultiPoint{{10, 20}, {30, 40}}
	loc.MultiPointPtr = &multiPointPtr

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Timestamps{},
			geom.Point{}, geom.LineString{}, geom.Polygon{},
			geom.MultiPoint{}, geom.MultiLineString{}, geom.MultiPolygon{},
			sfgeom.Point{}, sfgeom.LineString{}, sfgeom.Polygon{},
			sfgeom.MultiPoint{}, sfgeom.MultiLineString{}, sfgeom.MultiPolygon{},
		),
	)

	assert.Equal(t, len(loc.MultiPoint), len(out.MultiPoint))
	assert.Equal(t, len(loc.MultiLineString), len(out.MultiLineString))
	assert.Equal(t, len(loc.MultiPolygon), len(out.MultiPolygon))
}

func TestGeoCollection(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test Collection")

	loc.Collection = orb.Collection{
		orb.Point{1, 2},
		orb.LineString{{3, 4}, {5, 6}},
		orb.Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
	}

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Timestamps{},
			geom.Point{}, geom.LineString{}, geom.Polygon{},
			geom.MultiPoint{}, geom.MultiLineString{}, geom.MultiPolygon{},
			sfgeom.Point{}, sfgeom.LineString{}, sfgeom.Polygon{},
			sfgeom.MultiPoint{}, sfgeom.MultiLineString{}, sfgeom.MultiPolygon{},
		),
	)

	assert.Equal(t, len(loc.Collection), len(out.Collection))
}

func TestGeoNilPointers(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test Nil Pointers")
	loc.PointPtr = nil
	loc.LineStringPtr = nil
	loc.PolygonPtr = nil
	loc.MultiPointPtr = nil
	loc.GGPointPtr = nil
	loc.GGLineStringPtr = nil
	loc.GGPolygonPtr = nil
	loc.GGMultiPointPtr = nil
	loc.SFPointPtr = nil
	loc.SFLineStringPtr = nil
	loc.SFPolygonPtr = nil
	loc.SFMultiPointPtr = nil

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.Check(t, out.PointPtr == nil)
	assert.Check(t, out.LineStringPtr == nil)
	assert.Check(t, out.PolygonPtr == nil)
	assert.Check(t, out.MultiPointPtr == nil)
	assert.Check(t, out.GGPointPtr == nil)
	assert.Check(t, out.GGLineStringPtr == nil)
	assert.Check(t, out.GGPolygonPtr == nil)
	assert.Check(t, out.GGMultiPointPtr == nil)
	assert.Check(t, out.SFPointPtr == nil)
	assert.Check(t, out.SFLineStringPtr == nil)
	assert.Check(t, out.SFPolygonPtr == nil)
	assert.Check(t, out.SFMultiPointPtr == nil)
}

func TestGeoAllTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	pointPtr := orb.Point{100, 200}
	lineStringPtr := orb.LineString{{10, 20}, {30, 40}}
	polygonPtr := orb.Polygon{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}}
	multiPointPtr := orb.MultiPoint{{1, 1}, {2, 2}}

	ggPointPtr := geom.NewPointFlat(geom.XY, []float64{100, 200})
	ggLineStringPtr := geom.NewLineStringFlat(geom.XY, []float64{10, 20, 30, 40})
	ggPolygonPtr := geom.NewPolygonFlat(geom.XY, []float64{0, 0, 5, 0, 5, 5, 0, 5, 0, 0}, []int{5})
	ggMultiPointPtr := geom.NewMultiPointFlat(geom.XY, []float64{1, 1, 2, 2})

	sfPointPtr := sfNewPoint(100, 200)
	sfLineStringPtr := sfNewLineString([]float64{10, 20, 30, 40})
	sfPolygonPtr := sfNewPolygon([][]float64{{0, 0, 5, 0, 5, 5, 0, 5, 0, 0}})
	sfMultiPointPtr := sfgeom.NewMultiPoint([]sfgeom.Point{sfNewPoint(1, 1), sfNewPoint(2, 2)})

	loc := &model.Location{
		Name:            "Full Location Test",
		Point:           orb.Point{1.5, 2.5},
		PointPtr:        &pointPtr,
		LineString:      orb.LineString{{0, 0}, {1, 1}, {2, 2}},
		LineStringPtr:   &lineStringPtr,
		Polygon:         orb.Polygon{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
		PolygonPtr:      &polygonPtr,
		MultiPoint:      orb.MultiPoint{{0, 0}, {5, 5}, {10, 10}},
		MultiPointPtr:   &multiPointPtr,
		MultiLineString: orb.MultiLineString{{{0, 0}, {1, 1}}, {{2, 2}, {3, 3}}},
		MultiPolygon: orb.MultiPolygon{
			{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
		},
		Collection: orb.Collection{
			orb.Point{50, 50},
			orb.LineString{{60, 60}, {70, 70}},
		},

		GGPoint:           *geom.NewPointFlat(geom.XY, []float64{1.5, 2.5}),
		GGPointPtr:        ggPointPtr,
		GGLineString:      *geom.NewLineStringFlat(geom.XY, []float64{0, 0, 1, 1, 2, 2}),
		GGLineStringPtr:   ggLineStringPtr,
		GGPolygon:         *geom.NewPolygonFlat(geom.XY, []float64{0, 0, 10, 0, 10, 10, 0, 10, 0, 0}, []int{5}),
		GGPolygonPtr:      ggPolygonPtr,
		GGMultiPoint:      *geom.NewMultiPointFlat(geom.XY, []float64{0, 0, 5, 5, 10, 10}),
		GGMultiPointPtr:   ggMultiPointPtr,
		GGMultiLineString: *geom.NewMultiLineStringFlat(geom.XY, []float64{0, 0, 1, 1, 2, 2, 3, 3}, []int{2, 4}),
		GGMultiPolygon:    *geom.NewMultiPolygonFlat(geom.XY, []float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, [][]int{{5}}),

		SFPoint:           sfNewPoint(1.5, 2.5),
		SFPointPtr:        &sfPointPtr,
		SFLineString:      sfNewLineString([]float64{0, 0, 1, 1, 2, 2}),
		SFLineStringPtr:   &sfLineStringPtr,
		SFPolygon:         sfNewPolygon([][]float64{{0, 0, 10, 0, 10, 10, 0, 10, 0, 0}}),
		SFPolygonPtr:      &sfPolygonPtr,
		SFMultiPoint:      sfgeom.NewMultiPoint([]sfgeom.Point{sfNewPoint(0, 0), sfNewPoint(5, 5), sfNewPoint(10, 10)}),
		SFMultiPointPtr:   &sfMultiPointPtr,
		SFMultiLineString: sfgeom.NewMultiLineString([]sfgeom.LineString{sfNewLineString([]float64{0, 0, 1, 1}), sfNewLineString([]float64{2, 2, 3, 3})}),
		SFMultiPolygon:    sfgeom.NewMultiPolygon([]sfgeom.Polygon{sfNewPolygon([][]float64{{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}})}),
	}

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Timestamps{},
			geom.Point{}, geom.LineString{}, geom.Polygon{},
			geom.MultiPoint{}, geom.MultiLineString{}, geom.MultiPolygon{},
			sfgeom.Point{}, sfgeom.LineString{}, sfgeom.Polygon{},
			sfgeom.MultiPoint{}, sfgeom.MultiLineString{}, sfgeom.MultiPolygon{},
		),
	)
}

// go-geom specific tests

func TestGoGeomPoint(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test GoGeom Point")
	loc.GGPoint = *geom.NewPointFlat(geom.XY, []float64{1.0, 2.0})
	loc.GGPointPtr = geom.NewPointFlat(geom.XY, []float64{2.0, 3.0})

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.Equal(t, loc.GGPoint.X(), out.GGPoint.X())
	assert.Equal(t, loc.GGPoint.Y(), out.GGPoint.Y())
	assert.Check(t, out.GGPointPtr != nil)
	assert.Equal(t, loc.GGPointPtr.X(), out.GGPointPtr.X())
	assert.Equal(t, loc.GGPointPtr.Y(), out.GGPointPtr.Y())
}

func TestGoGeomLineString(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test GoGeom LineString")
	loc.GGLineString = *geom.NewLineStringFlat(geom.XY, []float64{1, 2, 3, 4, 5, 6})
	loc.GGLineStringPtr = geom.NewLineStringFlat(geom.XY, []float64{3, 4, 5, 6})

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.Equal(t, loc.GGLineString.NumCoords(), out.GGLineString.NumCoords())
	assert.DeepEqual(t, loc.GGLineString.FlatCoords(), out.GGLineString.FlatCoords())
	assert.Check(t, out.GGLineStringPtr != nil)
	assert.DeepEqual(t, loc.GGLineStringPtr.FlatCoords(), out.GGLineStringPtr.FlatCoords())
}

func TestGoGeomPolygon(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test GoGeom Polygon")
	loc.GGPolygon = *geom.NewPolygonFlat(geom.XY, []float64{0, 0, 10, 0, 10, 10, 0, 10, 0, 0}, []int{5})

	// Polygon with hole
	polyWithHole := geom.NewPolygonFlat(geom.XY,
		[]float64{0, 0, 10, 0, 10, 10, 0, 10, 0, 0, 2, 2, 8, 2, 8, 8, 2, 8, 2, 2},
		[]int{5, 10},
	)
	loc.GGPolygonPtr = polyWithHole

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.Equal(t, loc.GGPolygon.NumLinearRings(), out.GGPolygon.NumLinearRings())
	assert.DeepEqual(t, loc.GGPolygon.FlatCoords(), out.GGPolygon.FlatCoords())
	assert.Check(t, out.GGPolygonPtr != nil)
	assert.Equal(t, loc.GGPolygonPtr.NumLinearRings(), out.GGPolygonPtr.NumLinearRings())
}

func TestGoGeomMultiTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test GoGeom Multi Types")
	loc.GGMultiPoint = *geom.NewMultiPointFlat(geom.XY, []float64{1, 2, 3, 4, 5, 6})
	loc.GGMultiLineString = *geom.NewMultiLineStringFlat(geom.XY,
		[]float64{0, 0, 1, 1, 2, 2, 3, 3, 4, 4},
		[]int{2, 5},
	)
	loc.GGMultiPolygon = *geom.NewMultiPolygonFlat(geom.XY,
		[]float64{0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 2, 2, 3, 2, 3, 3, 2, 3, 2, 2},
		[][]int{{5}, {10}},
	)

	ggMultiPointPtr := geom.NewMultiPointFlat(geom.XY, []float64{10, 20, 30, 40})
	loc.GGMultiPointPtr = ggMultiPointPtr

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc.GGMultiPoint.FlatCoords(), out.GGMultiPoint.FlatCoords())
	assert.DeepEqual(t, loc.GGMultiLineString.FlatCoords(), out.GGMultiLineString.FlatCoords())
	assert.DeepEqual(t, loc.GGMultiPolygon.FlatCoords(), out.GGMultiPolygon.FlatCoords())
	assert.Check(t, out.GGMultiPointPtr != nil)
	assert.DeepEqual(t, loc.GGMultiPointPtr.FlatCoords(), out.GGMultiPointPtr.FlatCoords())
}

// simplefeatures specific tests

func TestSFPoint(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test SF Point")
	loc.SFPoint = sfNewPoint(1.0, 2.0)
	ptr := sfNewPoint(2.0, 3.0)
	loc.SFPointPtr = &ptr

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	inXY, _ := loc.SFPoint.XY()
	outXY, _ := out.SFPoint.XY()
	assert.Equal(t, inXY.X, outXY.X)
	assert.Equal(t, inXY.Y, outXY.Y)
	assert.Check(t, out.SFPointPtr != nil)
	inPtrXY, _ := loc.SFPointPtr.XY()
	outPtrXY, _ := out.SFPointPtr.XY()
	assert.Equal(t, inPtrXY.X, outPtrXY.X)
	assert.Equal(t, inPtrXY.Y, outPtrXY.Y)
}

func TestSFLineString(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test SF LineString")
	loc.SFLineString = sfNewLineString([]float64{1, 2, 3, 4, 5, 6})
	ptr := sfNewLineString([]float64{3, 4, 5, 6})
	loc.SFLineStringPtr = &ptr

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.Equal(t, loc.SFLineString.Coordinates().Length(), out.SFLineString.Coordinates().Length())
	assert.Check(t, out.SFLineStringPtr != nil)
	assert.Equal(t, loc.SFLineStringPtr.Coordinates().Length(), out.SFLineStringPtr.Coordinates().Length())
}

func TestSFPolygon(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test SF Polygon")
	loc.SFPolygon = sfNewPolygon([][]float64{{0, 0, 10, 0, 10, 10, 0, 10, 0, 0}})

	// Polygon with hole
	polyWithHole := sfNewPolygon([][]float64{
		{0, 0, 10, 0, 10, 10, 0, 10, 0, 0},
		{2, 2, 8, 2, 8, 8, 2, 8, 2, 2},
	})
	loc.SFPolygonPtr = &polyWithHole

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.Equal(t, loc.SFPolygon.ExteriorRing().Coordinates().Length(), out.SFPolygon.ExteriorRing().Coordinates().Length())
	assert.Check(t, out.SFPolygonPtr != nil)
	assert.Equal(t, loc.SFPolygonPtr.NumInteriorRings(), out.SFPolygonPtr.NumInteriorRings())
}

func TestSFMultiTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	loc := minimalLocation("Test SF Multi Types")
	loc.SFMultiPoint = sfgeom.NewMultiPoint([]sfgeom.Point{
		sfNewPoint(1, 2), sfNewPoint(3, 4), sfNewPoint(5, 6),
	})
	loc.SFMultiLineString = sfgeom.NewMultiLineString([]sfgeom.LineString{
		sfNewLineString([]float64{0, 0, 1, 1}),
		sfNewLineString([]float64{2, 2, 3, 3, 4, 4}),
	})
	loc.SFMultiPolygon = sfgeom.NewMultiPolygon([]sfgeom.Polygon{
		sfNewPolygon([][]float64{{0, 0, 1, 0, 1, 1, 0, 1, 0, 0}}),
		sfNewPolygon([][]float64{{2, 2, 3, 2, 3, 3, 2, 3, 2, 2}}),
	})

	sfMultiPointPtr := sfgeom.NewMultiPoint([]sfgeom.Point{sfNewPoint(10, 20), sfNewPoint(30, 40)})
	loc.SFMultiPointPtr = &sfMultiPointPtr

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, string(loc.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.Equal(t, loc.SFMultiPoint.NumPoints(), out.SFMultiPoint.NumPoints())
	assert.Equal(t, loc.SFMultiLineString.NumLineStrings(), out.SFMultiLineString.NumLineStrings())
	assert.Equal(t, loc.SFMultiPolygon.NumPolygons(), out.SFMultiPolygon.NumPolygons())
	assert.Check(t, out.SFMultiPointPtr != nil)
	assert.Equal(t, loc.SFMultiPointPtr.NumPoints(), out.SFMultiPointPtr.NumPoints())
}
