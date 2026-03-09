package basic

import (
	"context"
	"testing"
	"time"

	"som.test/gen/som/by"
	"som.test/gen/som/filter"
	"som.test/gen/som/repo"
	"som.test/model"
	"gotest.tools/v3/assert"
)

func TestBetweenDescribe(t *testing.T) {
	client := &repo.ClientImpl{}

	t.Run("default both inclusive", func(t *testing.T) {
		desc := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(10, 20)).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_int ∈ $A..=$B)",
			desc,
		)
	})

	t.Run("from exclusive", func(t *testing.T) {
		desc := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(10, 20).FromExclusive()).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_int ∈ $A>..=$B)",
			desc,
		)
	})

	t.Run("to exclusive", func(t *testing.T) {
		desc := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(10, 20).ToExclusive()).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_int ∈ $A..$B)",
			desc,
		)
	})

	t.Run("both exclusive", func(t *testing.T) {
		desc := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(10, 20).BothExclusive()).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_int ∈ $A>..$B)",
			desc,
		)
	})

	t.Run("float field", func(t *testing.T) {
		desc := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldFloat64.Between(1.5, 9.5)).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_float_64 ∈ $A..=$B)",
			desc,
		)
	})

	t.Run("string field", func(t *testing.T) {
		desc := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldString.Between("aaa", "zzz")).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_string ∈ $A..=$B)",
			desc,
		)
	})

	t.Run("time field", func(t *testing.T) {
		start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

		desc := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldTime.Between(start, end)).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_time ∈ $A..=$B)",
			desc,
		)
	})

	t.Run("combined with other filters", func(t *testing.T) {
		desc := client.AllTypesRepo().Query().
			Where(
				filter.AllTypes.FieldInt.Between(10, 20),
				filter.AllTypes.FieldString.Equal("hello"),
			).
			Describe()

		assert.Equal(t,
			"SELECT * FROM all_types WHERE (field_int ∈ $A..=$B AND field_string = $C)",
			desc,
		)
	})
}

func TestBetweenQuery(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 10; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldInt:   i * 10,
			FieldMonth: time.January,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("default both inclusive", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(20, 50)).
			Order(by.AllTypes.FieldInt.Asc()).
			All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 4, len(results))
		assert.Equal(t, 20, results[0].FieldInt)
		assert.Equal(t, 50, results[3].FieldInt)
	})

	t.Run("from exclusive", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(20, 50).FromExclusive()).
			Order(by.AllTypes.FieldInt.Asc()).
			All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 3, len(results))
		assert.Equal(t, 30, results[0].FieldInt)
		assert.Equal(t, 50, results[2].FieldInt)
	})

	t.Run("to exclusive", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(20, 50).ToExclusive()).
			Order(by.AllTypes.FieldInt.Asc()).
			All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 3, len(results))
		assert.Equal(t, 20, results[0].FieldInt)
		assert.Equal(t, 40, results[2].FieldInt)
	})

	t.Run("both exclusive", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(20, 50).BothExclusive()).
			Order(by.AllTypes.FieldInt.Asc()).
			All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, 30, results[0].FieldInt)
		assert.Equal(t, 40, results[1].FieldInt)
	})

	t.Run("empty result", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(91, 99)).
			All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 0, len(results))
	})

	t.Run("single boundary value", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Where(filter.AllTypes.FieldInt.Between(30, 30)).
			All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, 30, results[0].FieldInt)
	})
}
