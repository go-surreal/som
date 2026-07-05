package basic

import (
	"context"
	"testing"
	"time"

	"som.test/gen/som/by"
	"som.test/gen/som/query"
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
		Paginate(ctx, query.First(2), query.WithTotalCount())
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
		Paginate(ctx, query.First(2), query.After(page1.PageInfo.EndCursor))
	assert.NilError(t, err)

	assert.Equal(t, 2, len(page2.Items))
	assert.Equal(t, "c", page2.Items[0].FieldString)
	assert.Equal(t, "d", page2.Items[1].FieldString)
	assert.Equal(t, true, page2.PageInfo.HasNextPage)
	assert.Equal(t, true, page2.PageInfo.HasPreviousPage)

	// Page 3: last item, no further pages.
	page3, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate(ctx, query.First(2), query.After(page2.PageInfo.EndCursor))
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
		Paginate(ctx, query.Last(2))
	assert.NilError(t, err)

	assert.Equal(t, 2, len(last.Items))
	assert.Equal(t, "d", last.Items[0].FieldString)
	assert.Equal(t, "e", last.Items[1].FieldString)
	assert.Equal(t, true, last.PageInfo.HasPreviousPage)
	assert.Equal(t, false, last.PageInfo.HasNextPage)

	// Before the current page's start cursor → "b", "c".
	prev, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldString.Asc()).
		Paginate(ctx, query.Last(2), query.Before(last.PageInfo.StartCursor))
	assert.NilError(t, err)

	assert.Equal(t, 2, len(prev.Items))
	assert.Equal(t, "b", prev.Items[0].FieldString)
	assert.Equal(t, "c", prev.Items[1].FieldString)
	assert.Equal(t, true, prev.PageInfo.HasPreviousPage)
	assert.Equal(t, true, prev.PageInfo.HasNextPage)
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
		Paginate(ctx, query.First(2))
	assert.NilError(t, err)

	assert.Equal(t, 2, len(page1.Items))
	// Descending by time: newest (day 5) first, then day 4.
	assert.Equal(t, 4, page1.Items[0].FieldInt)
	assert.Equal(t, 3, page1.Items[1].FieldInt)
	assert.Equal(t, true, page1.PageInfo.HasNextPage)

	page2, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldTime.Desc()).
		Paginate(ctx, query.First(2), query.After(page1.PageInfo.EndCursor))
	assert.NilError(t, err)

	assert.Equal(t, 2, len(page2.Items))
	assert.Equal(t, 2, page2.Items[0].FieldInt)
	assert.Equal(t, 1, page2.Items[1].FieldInt)
}
