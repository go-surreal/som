package basic

import (
	"context"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/paulmach/orb"
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
	}
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

	out, exists, err := client.LocationRepo().Read(ctx, loc.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
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

	out, exists, err := client.LocationRepo().Read(ctx, loc.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
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

	out, exists, err := client.LocationRepo().Read(ctx, loc.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
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

	out, exists, err := client.LocationRepo().Read(ctx, loc.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
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

	out, exists, err := client.LocationRepo().Read(ctx, loc.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
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

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, loc.ID())
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
}

func TestGeoAllTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	pointPtr := orb.Point{100, 200}
	lineStringPtr := orb.LineString{{10, 20}, {30, 40}}
	polygonPtr := orb.Polygon{{{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0}}}
	multiPointPtr := orb.MultiPoint{{1, 1}, {2, 2}}

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
	}

	err := client.LocationRepo().Create(ctx, loc)
	if err != nil {
		t.Fatal(err)
	}

	out, exists, err := client.LocationRepo().Read(ctx, loc.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("location not found")
	}

	assert.DeepEqual(t, loc, out,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
	)
}
