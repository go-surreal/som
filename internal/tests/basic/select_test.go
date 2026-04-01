package basic

import (
	"context"
	"net/url"
	"sort"
	"testing"
	"time"

	som "som.test/gen/som"
	"som.test/gen/som/filter"
	"som.test/model"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestSelectAll(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("EmptyTable", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldString().All(ctx)
		assert.NilError(t, err)
		assert.Check(t, vals == nil)
	})

	now := time.Now().Truncate(time.Microsecond).UTC()
	time1 := now.Add(-2 * time.Hour)
	time2 := now.Add(-1 * time.Hour)
	time3 := now

	dur1 := time.Minute
	dur2 := time.Hour
	dur3 := 24 * time.Hour

	uuid1 := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	uuid2 := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	uuid3 := uuid.MustParse("00000000-0000-0000-0000-000000000003")

	url1, _ := url.Parse("https://example.com/alpha")
	url2, _ := url.Parse("https://example.com/bravo")
	url3, _ := url.Parse("https://example.com/charlie")

	records := []*model.AllTypes{
		{
			FieldString:      "alpha",
			FieldInt:         10,
			FieldFloat64:     1.1,
			FieldBool:        false,
			FieldTime:        time1,
			FieldDuration:    dur1,
			FieldMonth:       time.January,
			FieldWeekday:     time.Monday,
			FieldUUID:        uuid1,
			FieldURL:         *url1,
			FieldEmail:       som.Email("alice@example.com"),
			FieldEnum:        model.RoleAdmin,
			FieldCredentials: model.Credentials{Username: "alice", Password: "pass1"},
		},
		{
			FieldString:      "bravo",
			FieldInt:         20,
			FieldFloat64:     2.2,
			FieldBool:        true,
			FieldTime:        time2,
			FieldDuration:    dur2,
			FieldMonth:       time.February,
			FieldWeekday:     time.Tuesday,
			FieldUUID:        uuid2,
			FieldURL:         *url2,
			FieldEmail:       som.Email("bob@example.com"),
			FieldEnum:        model.RoleUser,
			FieldCredentials: model.Credentials{Username: "bob", Password: "pass2"},
		},
		{
			FieldString:      "charlie",
			FieldInt:         30,
			FieldFloat64:     3.3,
			FieldBool:        true,
			FieldTime:        time3,
			FieldDuration:    dur3,
			FieldMonth:       time.March,
			FieldWeekday:     time.Wednesday,
			FieldUUID:        uuid3,
			FieldURL:         *url3,
			FieldEmail:       som.Email("charlie@example.com"),
			FieldEnum:        model.RoleUser,
			FieldCredentials: model.Credentials{Username: "charlie", Password: "pass3"},
		},
	}

	for _, r := range records {
		if err := client.AllTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	// Primitive types

	t.Run("String", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldString().All(ctx)
		assert.NilError(t, err)
		sort.Strings(vals)
		assert.DeepEqual(t, []string{"alpha", "bravo", "charlie"}, vals)
	})

	t.Run("Int", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldInt().All(ctx)
		assert.NilError(t, err)
		sort.Ints(vals)
		assert.DeepEqual(t, []int{10, 20, 30}, vals)
	})

	t.Run("Float64", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldFloat64().All(ctx)
		assert.NilError(t, err)
		sort.Float64s(vals)
		assert.DeepEqual(t, []float64{1.1, 2.2, 3.3}, vals)
	})

	t.Run("Bool", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldBool().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return !vals[i] && vals[j] })
		assert.DeepEqual(t, []bool{false, true, true}, vals)
	})

	// Non-primitive types

	t.Run("Time", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldTime().All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 3, len(vals))
		sort.Slice(vals, func(i, j int) bool { return vals[i].Before(vals[j]) })
		assert.Check(t, vals[0].Equal(time1))
		assert.Check(t, vals[1].Equal(time2))
		assert.Check(t, vals[2].Equal(time3))
	})

	t.Run("Duration", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldDuration().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		assert.DeepEqual(t, []time.Duration{dur1, dur2, dur3}, vals)
	})

	t.Run("Month", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldMonth().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		assert.DeepEqual(t, []time.Month{time.January, time.February, time.March}, vals)
	})

	t.Run("Weekday", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldWeekday().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		assert.DeepEqual(t, []time.Weekday{time.Monday, time.Tuesday, time.Wednesday}, vals)
	})

	t.Run("UUID", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldUUID().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i].String() < vals[j].String() })
		assert.DeepEqual(t, []uuid.UUID{uuid1, uuid2, uuid3}, vals)
	})

	t.Run("URL", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldURL().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i].String() < vals[j].String() })
		assert.DeepEqual(t, []url.URL{*url1, *url2, *url3}, vals)
	})

	t.Run("Email", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldEmail().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		assert.DeepEqual(t, []som.Email{"alice@example.com", "bob@example.com", "charlie@example.com"}, vals)
	})

	t.Run("Enum", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldEnum().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		assert.DeepEqual(t, []model.Role{model.RoleAdmin, model.RoleUser, model.RoleUser}, vals)
	})

	// With filter

	t.Run("WithFilter", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.GreaterThan(15)).
			Select().FieldString().All(ctx)
		assert.NilError(t, err)
		sort.Strings(vals)
		assert.DeepEqual(t, []string{"bravo", "charlie"}, vals)
	})
}

func TestSelectFirst(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("EmptyTable", func(t *testing.T) {
		_, found, err := client.AllTypesRepo().Query().Select().FieldString().First(ctx)
		assert.NilError(t, err)
		assert.Check(t, !found)
	})

	records := []*model.AllTypes{
		{FieldString: "alpha", FieldInt: 10, FieldCredentials: model.Credentials{Username: "a", Password: "p"}, FieldMonth: 1},
		{FieldString: "bravo", FieldInt: 20, FieldCredentials: model.Credentials{Username: "b", Password: "p"}, FieldMonth: 1},
	}

	for _, r := range records {
		if err := client.AllTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("ReturnsOneValue", func(t *testing.T) {
		val, found, err := client.AllTypesRepo().Query().Select().FieldString().First(ctx)
		assert.NilError(t, err)
		assert.Check(t, found)
		assert.Check(t, val == "alpha" || val == "bravo")
	})

	t.Run("DoesNotCorruptAll", func(t *testing.T) {
		sf := client.AllTypesRepo().Query().Select().FieldString()

		// Call First, then All on the same SelectField — All should still return all records.
		_, _, err := sf.First(ctx)
		assert.NilError(t, err)

		vals, err := sf.All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 2, len(vals))
	})
}

func TestSelectDistinct(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	records := []*model.AllTypes{
		{FieldString: "alpha", FieldInt: 10, FieldCredentials: model.Credentials{Username: "a", Password: "p"}, FieldMonth: 1},
		{FieldString: "alpha", FieldInt: 20, FieldCredentials: model.Credentials{Username: "b", Password: "p"}, FieldMonth: 1},
		{FieldString: "bravo", FieldInt: 30, FieldCredentials: model.Credentials{Username: "c", Password: "p"}, FieldMonth: 1},
	}

	for _, r := range records {
		if err := client.AllTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("String", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldString().Distinct(ctx)
		assert.NilError(t, err)
		sort.Strings(vals)
		assert.DeepEqual(t, []string{"alpha", "bravo"}, vals)
	})

	t.Run("Int", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldInt().Distinct(ctx)
		assert.NilError(t, err)
		sort.Ints(vals)
		assert.DeepEqual(t, []int{10, 20, 30}, vals)
	})

	t.Run("WithFilter", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.GreaterThan(15)).
			Select().FieldString().Distinct(ctx)
		assert.NilError(t, err)
		sort.Strings(vals)
		assert.DeepEqual(t, []string{"alpha", "bravo"}, vals)
	})
}
