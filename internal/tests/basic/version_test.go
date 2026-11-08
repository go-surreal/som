package basic

import (
	"context"
	"reflect"
	"testing"
	"time"

	"som.test/gen/som/filter"
	"som.test/gen/som/query"
	"som.test/gen/som/repo"
	"som.test/gen/som/with"
	"som.test/model"
	"gotest.tools/v3/assert"
)

// The test container runs SurrealKV with versioning enabled and the test
// module is generated with --versioned, so Version() is available here.

func TestVersionReadsPastState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	author := model.SpecialTypes{Name: "before"}
	err := client.SpecialTypesRepo().Create(ctx, &author)
	assert.NilError(t, err)

	post := model.SpecialRelation{Title: "before", Author: &author}
	err = client.SpecialRelationRepo().Create(ctx, &post)
	assert.NilError(t, err)

	snapshot := serverNow(ctx, t, client)

	author.Name = "after"
	err = client.SpecialTypesRepo().Update(ctx, &author)
	assert.NilError(t, err)

	post.Title = "after"
	err = client.SpecialRelationRepo().Update(ctx, &post)
	assert.NilError(t, err)

	err = client.SpecialRelationRepo().Create(ctx, &model.SpecialRelation{Title: "later"})
	assert.NilError(t, err)

	past, err := client.SpecialRelationRepo().Query().
		Fetch(with.SpecialRelation.Author()).
		Version(snapshot).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(past), 1)
	assert.Equal(t, past[0].Title, "before")
	assert.Assert(t, past[0].Author != nil)
	assert.Equal(t, past[0].Author.Name, "before", "fetched links resolve at the same version")

	count, err := client.SpecialRelationRepo().Query().Version(snapshot).Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, count, 1)

	filtered, err := client.SpecialRelationRepo().Query().
		Where(filter.SpecialRelation.Title.Equal("after")).
		Version(snapshot).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(filtered), 0, "filters apply to the past state")

	current, err := client.SpecialRelationRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(current), 2)
}

// TestVersionTerminals verifies that the version also applies to the query
// terminals that issue more than one statement or seek by record ID range.
func TestVersionTerminals(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	records := make([]*model.AllTypes, 0, 3)

	for range 3 {
		record := &model.AllTypes{FieldString: "quick fox", FieldMonth: time.January}
		err := client.AllTypesRepo().Create(ctx, record)
		assert.NilError(t, err)

		records = append(records, record)
	}

	snapshot := serverNow(ctx, t, client)

	for _, record := range records {
		record.FieldString = "slow dog"
		err := client.AllTypesRepo().Update(ctx, record)
		assert.NilError(t, err)
	}

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{FieldString: "quick fox", FieldMonth: time.January})
	assert.NilError(t, err)

	query := client.AllTypesRepo().Query().Version(snapshot)

	var iterated []string
	for record, err := range query.Iterate(ctx, 2) {
		assert.NilError(t, err)
		iterated = append(iterated, record.FieldString)
	}
	assert.DeepEqual(t, iterated, []string{"quick fox", "quick fox", "quick fox"})

	firstPage, err := query.Paginate().First(2).WithTotalCount().Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(firstPage.Items), 2)
	assert.Equal(t, firstPage.TotalCount, 3)

	secondPage, err := query.Paginate().First(2).After(firstPage.NextCursor()).Get(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(secondPage.Items), 1)
	assert.Equal(t, secondPage.Items[0].FieldString, "quick fox")

}

// TestVersionCurrentStateOnly verifies that a versioned query does not offer
// the methods that only read the current state, nor the chain methods that
// would hand them back.
func TestVersionCurrentStateOnly(t *testing.T) {
	t.Parallel()

	versioned := reflect.TypeFor[query.Versioned[model.AllTypes]]()

	for _, method := range []string{"Live", "LiveCount", "Search", "SearchAll", "Where", "Order", "Fetch"} {
		_, hasMethod := versioned.MethodByName(method)
		assert.Assert(t, !hasMethod, "versioned query must not offer %s", method)
	}

	for _, method := range []string{"All", "Count", "Iterate", "Paginate", "Describe"} {
		_, hasMethod := versioned.MethodByName(method)
		assert.Assert(t, hasMethod, "versioned query must offer %s", method)
	}
}

func TestVersionBeforeTableExists(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	_, err := client.SpecialRelationRepo().Query().
		Version(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)).
		All(ctx)
	assert.ErrorContains(t, err, "does not exist")
}

// TestVersionExpiry verifies that expiry is evaluated relative to the version
// time, so a record that has expired since is still part of the past state.
func TestVersionExpiry(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.EphemeralRepo().Create(ctx, &model.Ephemeral{Label: "temp"})
	assert.NilError(t, err)

	snapshot := serverNow(ctx, t, client)

	time.Sleep(1500 * time.Millisecond)

	current, err := client.EphemeralRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(current), 0, "record has expired by now")

	past, err := client.EphemeralRepo().Query().Version(snapshot).All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(past), 1, "record had not expired at the version time")
}

// serverNow returns the current time of the database server. The versions are
// stamped by the server clock, which may drift from the one of the test host.
func serverNow(ctx context.Context, tb testing.TB, client repo.Client) time.Time {
	tb.Helper()

	result, err := client.Raw(ctx, "RETURN [<string> time::now()]", nil)
	assert.NilError(tb, err)

	formatted, err := result.ScanOne[string]()
	assert.NilError(tb, err)

	now, err := time.Parse(time.RFC3339Nano, formatted)
	assert.NilError(tb, err)

	return now
}
