package basic

import (
	"context"
	"math"
	"testing"

	"gotest.tools/v3/assert"
	"som.test/gen/som/filter"
	"som.test/gen/som/with"
	"som.test/model"
)

func TestUnionLinkRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	square := &model.Square{Label: "big", Size: 4}
	assert.NilError(t, client.SquareRepo().Create(ctx, square))

	circle := &model.Circle{Label: "round", Radius: 2}
	assert.NilError(t, client.CircleRepo().Create(ctx, circle))

	squared := &model.Drawing{Name: "squared", Subject: square}
	assert.NilError(t, client.DrawingRepo().Create(ctx, squared))

	rounded := &model.Drawing{Name: "rounded", Subject: circle}
	assert.NilError(t, client.DrawingRepo().Create(ctx, rounded))

	// An unfetched link holds the record id of its member only.
	read, ok, err := client.DrawingRepo().Read(ctx, string(squared.ID()))
	assert.NilError(t, err)
	assert.Assert(t, ok)

	link, isSquare := read.Subject.(*model.Square)
	assert.Assert(t, isSquare, "expected the subject to be a square, got %T", read.Subject)
	assert.Equal(t, link.ID(), square.ID())
	assert.Assert(t, link.IsPartial())

	// A fetched link holds the whole record of its member.
	fetched, err := client.DrawingRepo().Query().
		Where(filter.Drawing.Name.Equal("rounded")).
		Fetch(with.Drawing.Subject()).
		First(ctx)
	assert.NilError(t, err)

	loaded, isCircle := fetched.Subject.(*model.Circle)
	assert.Assert(t, isCircle, "expected the subject to be a circle, got %T", fetched.Subject)
	assert.Equal(t, loaded.Radius, 2.0)
	assert.Assert(t, !loaded.IsPartial())
	assert.Equal(t, fetched.Subject.Area(), math.Pi*4)
}

func TestUnionLinkFilter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	small := &model.Square{Label: "small", Size: 1}
	assert.NilError(t, client.SquareRepo().Create(ctx, small))

	big := &model.Square{Label: "big", Size: 9}
	assert.NilError(t, client.SquareRepo().Create(ctx, big))

	circle := &model.Circle{Label: "round", Radius: 2}
	assert.NilError(t, client.CircleRepo().Create(ctx, circle))

	for _, subject := range []model.Shape{small, big, circle} {
		assert.NilError(t, client.DrawingRepo().Create(ctx, &model.Drawing{
			Name:    "drawing",
			Subject: subject,
		}))
	}

	assert.NilError(t, client.DrawingRepo().Create(ctx, &model.Drawing{Name: "empty"}))

	// Match by the member table the link points to.
	squares, err := client.DrawingRepo().Query().
		Where(filter.Drawing.Subject().IsSquare()).
		Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, squares, 2)

	circles, err := client.DrawingRepo().Query().
		Where(filter.Drawing.Subject().IsCircle()).
		Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, circles, 1)

	// Narrow the link to a member and filter on that member's field.
	large, err := client.DrawingRepo().Query().
		Where(filter.Drawing.Subject().Square().Size.GreaterThan(3)).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(large), 1)

	unset, err := client.DrawingRepo().Query().
		Where(filter.Drawing.Subject().Nil(true)).
		Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, unset, 1)
}

func TestOpenUnionMembership(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	square := &model.Square{Label: "big", Size: 4}
	assert.NilError(t, client.SquareRepo().Create(ctx, square))

	circle := &model.Circle{Label: "round", Radius: 2}
	assert.NilError(t, client.CircleRepo().Create(ctx, circle))

	// Square and Circle are members of the sealed union Shape and of the open
	// union Colored at the same time.
	assert.NilError(t, client.DrawingRepo().Create(ctx, &model.Drawing{
		Name:    "both",
		Subject: square,
		Accent:  circle,
	}))

	read, err := client.DrawingRepo().Query().
		Where(filter.Drawing.Name.Equal("both")).
		Fetch(with.Drawing.Accent()).
		First(ctx)
	assert.NilError(t, err)

	accent, isCircle := read.Accent.(*model.Circle)
	assert.Assert(t, isCircle, "expected the accent to be a circle, got %T", read.Accent)
	assert.Equal(t, accent.Radius, 2.0)
	assert.Equal(t, read.Accent.Color(), "blue")

	squares, err := client.DrawingRepo().Query().
		Where(filter.Drawing.Accent().IsCircle()).
		Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, squares, 1)

	// The open union has a repository of its own, just like the sealed one.
	count, err := client.ColoredRepo().Query().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, count, 2)
}

func TestUnionRepoQuery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	assert.NilError(t, client.SquareRepo().Create(ctx, &model.Square{Label: "a", Size: 1}))
	assert.NilError(t, client.SquareRepo().Create(ctx, &model.Square{Label: "b", Size: 2}))
	assert.NilError(t, client.CircleRepo().Create(ctx, &model.Circle{Label: "c", Radius: 3}))

	count, err := client.ShapeRepo().Query().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, count, 3)

	shapes, err := client.ShapeRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(shapes), 3)

	var squares, circles int
	var area float64

	for _, shape := range shapes {
		// A method declared by the union interface is implemented by every
		// member, so it can be called without knowing the concrete one.
		area += shape.Area()

		switch shape.(type) {
		case *model.Square:
			squares++
		case *model.Circle:
			circles++
		default:
			t.Fatalf("unexpected member %T", shape)
		}
	}

	assert.Equal(t, squares, 2)
	assert.Equal(t, circles, 1)
	assert.Equal(t, area, 1+4+math.Pi*9)

	// A query over the union may narrow to a single member as well.
	narrowed, err := client.ShapeRepo().Query().
		Where(filter.Shape.Square().Size.GreaterThan(1)).
		Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, narrowed, 1)
}
