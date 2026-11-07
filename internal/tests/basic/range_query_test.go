package basic

import (
	"context"
	"testing"
	"time"

	som "som.test/gen/som"
	"som.test/gen/som/by"
	"som.test/gen/som/filter"
	"som.test/gen/som/query"
	"som.test/model"
	"gotest.tools/v3/assert"
)

func TestRangeQueryDescribe(t *testing.T) {
	t.Parallel()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	t.Run("ArrayID full range", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.From(model.WeatherKey{City: "London", Date: start}),
			som.To(model.WeatherKey{City: "London", Date: end}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('weather', [$A, $B]..[$C, $D])",
			desc,
		)
	})

	t.Run("ArrayID open-ended to end", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.From(model.WeatherKey{City: "London", Date: start}),
			som.ToEnd(),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('weather', [$A, $B]..)",
			desc,
		)
	})

	t.Run("ArrayID open-ended from start", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.FromStart(),
			som.To(model.WeatherKey{City: "London", Date: end}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('weather', ..[$A, $B])",
			desc,
		)
	})

	t.Run("ArrayID exclusive from inclusive to", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.FromExclusive(model.WeatherKey{City: "London", Date: start}),
			som.ToInclusive(model.WeatherKey{City: "London", Date: end}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('weather', [$A, $B]>..=[$C, $D])",
			desc,
		)
	})

	t.Run("ObjectID full range", func(t *testing.T) {
		desc := query.NewPersonObj(nil).Range(
			som.From(model.PersonKey{Name: "Alice", Age: 20}),
			som.To(model.PersonKey{Name: "Bob", Age: 30}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('person_obj', {name: $A, age: $B}..{name: $C, age: $D})",
			desc,
		)
	})

	t.Run("Range with WHERE", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.From(model.WeatherKey{City: "London", Date: start}),
			som.To(model.WeatherKey{City: "London", Date: end}),
		).Where(
			filter.Weather.Temperature.GreaterThan(20),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('weather', [$A, $B]..[$C, $D]) WHERE (temperature > $E)",
			desc,
		)
	})
}

func TestRangeQueryStringIDDescribe(t *testing.T) {
	t.Parallel()

	t.Run("ULID full range", func(t *testing.T) {
		desc := query.NewAllTypes(nil).Range(
			som.From(som.ULID("01HY5E8ZQA1BCD2EF3GH4JK5MN")),
			som.To(som.ULID("01HY5E8ZQA9ZZZ9ZZ9ZZ9ZZ9ZZ")),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('all_types', $A..$B)",
			desc,
		)
	})

	t.Run("ULID open-ended to end", func(t *testing.T) {
		desc := query.NewAllTypes(nil).Range(
			som.From(som.ULID("01HY5E8ZQA1BCD2EF3GH4JK5MN")),
			som.ToEnd(),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('all_types', $A..)",
			desc,
		)
	})

	t.Run("ULID open-ended from start", func(t *testing.T) {
		desc := query.NewAllTypes(nil).Range(
			som.FromStart(),
			som.To(som.ULID("01HY5E8ZQA9ZZZ9ZZ9ZZ9ZZ9ZZ")),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('all_types', ..$A)",
			desc,
		)
	})

	t.Run("ULID exclusive from inclusive to", func(t *testing.T) {
		desc := query.NewAllTypes(nil).Range(
			som.FromExclusive(som.ULID("01HY5E8ZQA1BCD2EF3GH4JK5MN")),
			som.ToInclusive(som.ULID("01HY5E8ZQA9ZZZ9ZZ9ZZ9ZZ9ZZ")),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('all_types', $A>..=$B)",
			desc,
		)
	})

	t.Run("ULID sole id DESC sort repeats the key", func(t *testing.T) {
		desc := query.NewAllTypes(nil).Range(
			som.FromStart(),
			som.ToEnd(),
		).Order(
			by.AllTypes.ID.Desc(),
		).Limit(2).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('all_types', ..) ORDER BY id DESC, id DESC LIMIT 2",
			desc,
		)
	})

	t.Run("ULID id DESC with further sorts is kept as is", func(t *testing.T) {
		desc := query.NewAllTypes(nil).Range(
			som.FromStart(),
			som.ToEnd(),
		).Order(
			by.AllTypes.ID.Desc(),
			by.AllTypes.FieldString.Asc(),
		).Limit(2).Describe()

		assert.Equal(t,
			"SELECT * FROM type::record('all_types', ..) ORDER BY id DESC, field_string ASC LIMIT 2",
			desc,
		)
	})
}

func TestRangeQueryStringID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ids := []som.String{"a", "b", "c", "d", "e", "f"}

	for _, id := range ids {
		err := client.SlugRepo().CreateWithID(ctx, string(id), &model.Slug{Title: "slug " + string(id)})
		assert.NilError(t, err)
	}

	resultIDs := func(records []*model.Slug) []som.String {
		var out []som.String

		for _, record := range records {
			out = append(out, record.ID())
		}

		return out
	}

	t.Run("bounded range", func(t *testing.T) {
		results, err := client.SlugRepo().Query().Range(
			som.From(ids[1]),
			som.ToInclusive(ids[4]),
		).Order(
			by.Slug.ID.Asc(),
		).All(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, ids[1:5], resultIDs(results))
	})

	t.Run("exclusive from", func(t *testing.T) {
		results, err := client.SlugRepo().Query().Range(
			som.FromExclusive(ids[3]),
			som.ToEnd(),
		).Order(
			by.Slug.ID.Asc(),
		).All(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, ids[4:], resultIDs(results))
	})

	t.Run("id ASC with limit", func(t *testing.T) {
		results, err := client.SlugRepo().Query().Range(
			som.FromStart(),
			som.ToEnd(),
		).Order(
			by.Slug.ID.Asc(),
		).Limit(2).All(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, ids[:2], resultIDs(results))
	})

	t.Run("id DESC with limit returns the highest ids", func(t *testing.T) {
		results, err := client.SlugRepo().Query().Range(
			som.FromStart(),
			som.To(ids[5]),
		).Order(
			by.Slug.ID.Desc(),
		).Limit(2).All(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, []som.String{ids[4], ids[3]}, resultIDs(results))
	})
}

func TestRangeQueryArrayID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	dates := []time.Time{
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC),
	}
	temps := []float64{5.0, 12.0, 25.0, 18.0, 3.0}

	for i, date := range dates {
		w := &model.Weather{
			Node:        som.NewNode[model.WeatherKey](model.WeatherKey{City: "London", Date: date}),
			Temperature: temps[i],
		}
		err := client.WeatherRepo().CreateWithID(ctx, w)
		assert.NilError(t, err)
	}

	t.Run("full range", func(t *testing.T) {
		from := model.WeatherKey{City: "London", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}
		to := model.WeatherKey{City: "London", Date: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)}

		results, err := client.WeatherRepo().Query().Range(
			som.From(from),
			som.To(to),
		).All(ctx)
		assert.NilError(t, err)
		// March 15, June 20, Sept 25 are in [March 1, Oct 1)
		assert.Equal(t, 3, len(results))
	})

	t.Run("open-ended to end", func(t *testing.T) {
		from := model.WeatherKey{City: "London", Date: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)}

		results, err := client.WeatherRepo().Query().Range(
			som.From(from),
			som.ToEnd(),
		).All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 3, len(results))
	})

	t.Run("open-ended from start", func(t *testing.T) {
		to := model.WeatherKey{City: "London", Date: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)}

		results, err := client.WeatherRepo().Query().Range(
			som.FromStart(),
			som.To(to),
		).All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 2, len(results))
	})

	t.Run("range with filter", func(t *testing.T) {
		from := model.WeatherKey{City: "London", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
		to := model.WeatherKey{City: "London", Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}

		results, err := client.WeatherRepo().Query().Range(
			som.From(from),
			som.To(to),
		).Where(
			filter.Weather.Temperature.GreaterThan(10),
		).All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 3, len(results))
		for _, r := range results {
			assert.Assert(t, r.Temperature > 10)
		}
	})

	t.Run("range with order", func(t *testing.T) {
		from := model.WeatherKey{City: "London", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
		to := model.WeatherKey{City: "London", Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}

		results, err := client.WeatherRepo().Query().Range(
			som.From(from),
			som.To(to),
		).Order(
			by.Weather.Temperature.Desc(),
		).All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 5, len(results))
		assert.Equal(t, 25.0, results[0].Temperature)
		assert.Equal(t, 3.0, results[len(results)-1].Temperature)
	})

	t.Run("range with count", func(t *testing.T) {
		from := model.WeatherKey{City: "London", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
		to := model.WeatherKey{City: "London", Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}

		count, err := client.WeatherRepo().Query().Range(
			som.From(from),
			som.To(to),
		).Count(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("empty range", func(t *testing.T) {
		from := model.WeatherKey{City: "London", Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
		to := model.WeatherKey{City: "London", Date: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)}

		results, err := client.WeatherRepo().Query().Range(
			som.From(from),
			som.To(to),
		).All(ctx)
		assert.NilError(t, err)
		assert.Equal(t, 0, len(results))
	})
}

func TestRangeQueryObjectID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	people := []struct {
		key   model.PersonKey
		email string
	}{
		{model.PersonKey{Name: "Alice", Age: 25}, "alice@example.com"},
		{model.PersonKey{Name: "Bob", Age: 30}, "bob@example.com"},
		{model.PersonKey{Name: "Charlie", Age: 35}, "charlie@example.com"},
		{model.PersonKey{Name: "Diana", Age: 40}, "diana@example.com"},
	}

	for _, p := range people {
		person := &model.PersonObj{
			Node:  som.NewNode[model.PersonKey](p.key),
			Email: p.email,
		}
		err := client.PersonObjRepo().CreateWithID(ctx, person)
		assert.NilError(t, err)
	}

	t.Run("full range", func(t *testing.T) {
		results, err := client.PersonObjRepo().Query().Range(
			som.From(model.PersonKey{Name: "Bob", Age: 30}),
			som.To(model.PersonKey{Name: "Diana", Age: 40}),
		).All(ctx)
		assert.NilError(t, err)
		for _, r := range results {
			t.Logf("got: %s (age %d)", r.ID().Name, r.ID().Age)
		}
		assert.Assert(t, len(results) > 0)
	})

	t.Run("inclusive to", func(t *testing.T) {
		results, err := client.PersonObjRepo().Query().Range(
			som.From(model.PersonKey{Name: "Bob", Age: 30}),
			som.ToInclusive(model.PersonKey{Name: "Diana", Age: 40}),
		).All(ctx)
		assert.NilError(t, err)
		for _, r := range results {
			t.Logf("got: %s (age %d)", r.ID().Name, r.ID().Age)
		}
		assert.Assert(t, len(results) > 0)
	})
}
