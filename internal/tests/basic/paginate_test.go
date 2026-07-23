package basic

import (
	"context"
	"strings"
	"testing"
	"time"

	"som.test/gen/som"
	"som.test/gen/som/by"
	"som.test/model"
	"gotest.tools/v3/assert"
)

func TestPaginateForward(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	names := []string{"a", "b", "c", "d", "e"}
	for i, name := range names {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: name,
			FieldInt:    i,
			FieldTime:   time.Date(2020, time.January, i+1, 0, 0, 0, 0, time.UTC),
			FieldMonth:  time.January,
		})
		assert.NilError(t, err)
	}

	// Page 1: first 2 ordered by FieldString ascending.
	page1, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().First(2).WithTotalCount().Get(ctx)
	assert.NilError(t, err)

	assert.Equal(t, 2, len(page1.Items))
	assert.Equal(t, "a", page1.Items[0].FieldString)
	assert.Equal(t, "b", page1.Items[1].FieldString)
	assert.Equal(t, 5, page1.TotalCount)
	assert.Equal(t, true, page1.PageInfo.HasNextPage)
	assert.Equal(t, false, page1.PageInfo.HasPreviousPage)
	assert.Assert(t, page1.PageInfo.EndCursor != "")

	// Page 2: continue after page 1's end cursor.
	page2, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().First(2).After(page1.PageInfo.EndCursor).Get(ctx)
	assert.NilError(t, err)

	assert.Equal(t, 2, len(page2.Items))
	assert.Equal(t, "c", page2.Items[0].FieldString)
	assert.Equal(t, "d", page2.Items[1].FieldString)
	assert.Equal(t, true, page2.PageInfo.HasNextPage)
	assert.Equal(t, true, page2.PageInfo.HasPreviousPage)

	// Page 3: last item, no further pages.
	page3, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().First(2).After(page2.PageInfo.EndCursor).Get(ctx)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(page3.Items))
	assert.Equal(t, "e", page3.Items[0].FieldString)
	assert.Equal(t, false, page3.PageInfo.HasNextPage)
	assert.Equal(t, true, page3.PageInfo.HasPreviousPage)
}

func TestPaginateBackward(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	names := []string{"a", "b", "c", "d", "e"}
	for i, name := range names {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: name,
			FieldInt:    i,
			FieldMonth:  time.January,
		})
		assert.NilError(t, err)
	}

	// Last 2 ordered by FieldString ascending → "d", "e".
	last, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().Last(2).Get(ctx)
	assert.NilError(t, err)

	assert.Equal(t, 2, len(last.Items))
	assert.Equal(t, "d", last.Items[0].FieldString)
	assert.Equal(t, "e", last.Items[1].FieldString)
	assert.Equal(t, true, last.PageInfo.HasPreviousPage)
	assert.Equal(t, false, last.PageInfo.HasNextPage)

	// Before the current page's start cursor → "b", "c".
	prev, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().Last(2).Before(last.PageInfo.StartCursor).Get(ctx)
	assert.NilError(t, err)

	assert.Equal(t, 2, len(prev.Items))
	assert.Equal(t, "b", prev.Items[0].FieldString)
	assert.Equal(t, "c", prev.Items[1].FieldString)
	assert.Equal(t, true, prev.PageInfo.HasPreviousPage)
	assert.Equal(t, true, prev.PageInfo.HasNextPage)
}

// TestPaginateNextPrev exercises the Page.Next()/Prev() navigation helpers,
// which carry the original query, page size and options forward.
func TestPaginateNextPrev(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for _, n := range []string{"a", "b", "c", "d", "e"} {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: n,
			FieldMonth:  time.January,
		})
		assert.NilError(t, err)
	}

	page1, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().First(2).Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "a", page1.Items[0].FieldString)
	assert.Equal(t, "b", page1.Items[1].FieldString)

	page2, err := page1.Next().Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "c", page2.Items[0].FieldString)
	assert.Equal(t, "d", page2.Items[1].FieldString)

	// Back to page 1 from page 2.
	back, err := page2.Prev().Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(back.Items))
	assert.Equal(t, "a", back.Items[0].FieldString)
	assert.Equal(t, "b", back.Items[1].FieldString)
}

// TestPaginateTypedCursor exercises a non-string sort key (time.Time) to
// ensure cursor values are encoded with their DB type and compared correctly.
func TestPaginateTypedCursor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: string(rune('a' + i)),
			FieldInt:    i,
			FieldTime:   time.Date(2020, time.January, i+1, 12, 0, 0, 0, time.UTC),
			FieldMonth:  time.January,
		})
		assert.NilError(t, err)
	}

	page1, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldTime.Desc()).
		Paginate().First(2).Get(ctx)
	assert.NilError(t, err)

	assert.Equal(t, 2, len(page1.Items))
	// Descending by time: newest (day 5) first, then day 4.
	assert.Equal(t, 4, page1.Items[0].FieldInt)
	assert.Equal(t, 3, page1.Items[1].FieldInt)
	assert.Equal(t, true, page1.PageInfo.HasNextPage)

	page2, err := page1.Next().Get(ctx)
	assert.NilError(t, err)

	assert.Equal(t, 2, len(page2.Items))
	assert.Equal(t, 2, page2.Items[0].FieldInt)
	assert.Equal(t, 1, page2.Items[1].FieldInt)
}

// TestPaginateAccuratePageInfo exercises the WithAccuratePageInfo path, which
// runs an extra boundary query. It verifies the flags stay correct and, just
// as importantly, that the extra query is well-formed (does not error).
func TestPaginateAccuratePageInfo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for _, n := range []string{"a", "b", "c", "d", "e"} {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: n,
			FieldMonth:  time.January,
		})
		assert.NilError(t, err)
	}

	// First page, no cursor: there genuinely is no previous page.
	first, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().First(2).WithAccuratePageInfo().Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "a", first.Items[0].FieldString)
	assert.Equal(t, false, first.PageInfo.HasPreviousPage)
	assert.Equal(t, true, first.PageInfo.HasNextPage)

	// Middle page via cursor still reports a previous page.
	mid, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().First(2).After(first.PageInfo.EndCursor).WithAccuratePageInfo().Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "c", mid.Items[0].FieldString)
	assert.Equal(t, true, mid.PageInfo.HasPreviousPage)
	assert.Equal(t, true, mid.PageInfo.HasNextPage)

	// Last page (backward, no cursor): there genuinely is no next page.
	last, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate().Last(2).WithAccuratePageInfo().Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "e", last.Items[1].FieldString)
	assert.Equal(t, false, last.PageInfo.HasNextPage)
	assert.Equal(t, true, last.PageInfo.HasPreviousPage)
}

// TestPaginateComplexIDGuard verifies that pagination fails with a clear error
// for models with a complex ID, which have no single "id" tiebreaker and
// should use Range() instead.
func TestPaginateComplexIDGuard(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	person := &model.PersonObj{
		Node:  som.NewNode[model.PersonKey](model.PersonKey{Name: "Alice", Age: 30}),
		Email: "alice@example.com",
	}
	assert.NilError(t, client.PersonObjRepo().CreateWithID(ctx, person))

	_, err := client.PersonObjRepo().Query().Paginate().First(10).Get(ctx)
	assert.Assert(t, err != nil, "expected an error for complex-ID pagination")
	assert.Assert(t, strings.Contains(err.Error(), "Range()"),
		"error should point to Range(), got: %v", err)
}
