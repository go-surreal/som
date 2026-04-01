package basic

import (
	"context"
	"sort"
	"testing"

	"som.test/gen/som/filter"
	"som.test/model"
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

	records := []*model.AllTypes{
		{FieldString: "alpha", FieldInt: 10, FieldBool: false, FieldCredentials: model.Credentials{Username: "alice", Password: "pass1"}, FieldMonth: 1},
		{FieldString: "bravo", FieldInt: 20, FieldBool: true, FieldCredentials: model.Credentials{Username: "bob", Password: "pass2"}, FieldMonth: 1},
		{FieldString: "charlie", FieldInt: 30, FieldBool: true, FieldCredentials: model.Credentials{Username: "charlie", Password: "pass3"}, FieldMonth: 1},
	}

	for _, r := range records {
		if err := client.AllTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

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

	t.Run("Bool", func(t *testing.T) {
		vals, err := client.AllTypesRepo().Query().Select().FieldBool().All(ctx)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return !vals[i] && vals[j] })
		assert.DeepEqual(t, []bool{false, true, true}, vals)
	})

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
