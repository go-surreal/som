package basic

import (
	"cmp"
	"context"
	"slices"
	"testing"
	"time"

	som "som.test/gen/som"
	"som.test/gen/som/by"
	"som.test/gen/som/query"
	"som.test/gen/som/repo"
	"som.test/model"
	"gotest.tools/v3/assert"
)

type weatherReading struct {
	key         model.WeatherKey
	temperature float64
}

// seedWeather creates readings for two cities and returns them in record ID
// order (city first, then date).
func seedWeather(ctx context.Context, t *testing.T, client repo.Client) []weatherReading {
	t.Helper()

	day := func(dayOfMonth int) time.Time {
		return time.Date(2024, time.March, dayOfMonth, 0, 0, 0, 0, time.UTC)
	}

	readings := []weatherReading{
		{model.WeatherKey{City: "Berlin", Date: day(1)}, 12},
		{model.WeatherKey{City: "Berlin", Date: day(2)}, 18},
		{model.WeatherKey{City: "London", Date: day(1)}, 9},
		{model.WeatherKey{City: "London", Date: day(2)}, 18},
		{model.WeatherKey{City: "London", Date: day(3)}, 14},
		{model.WeatherKey{City: "London", Date: day(4)}, 18},
		{model.WeatherKey{City: "London", Date: day(5)}, 7},
		{model.WeatherKey{City: "London", Date: day(6)}, 11},
	}

	for _, reading := range readings {
		weather := &model.Weather{
			Node:        som.NewNode[model.WeatherKey](reading.key),
			Temperature: reading.temperature,
		}

		assert.NilError(t, client.WeatherRepo().CreateWithID(ctx, weather))
	}

	return readings
}

func readingKeys(readings []weatherReading) []model.WeatherKey {
	keys := make([]model.WeatherKey, len(readings))

	for index, reading := range readings {
		keys[index] = reading.key
	}

	return keys
}

func weatherKeys(items []*model.Weather) []model.WeatherKey {
	keys := make([]model.WeatherKey, len(items))

	for index, item := range items {
		keys[index] = item.ID()
	}

	return keys
}

// walkWeatherPages pages forward through a query until there is no next page,
// returning the keys of all items in page order.
func walkWeatherPages(ctx context.Context, t *testing.T, paginate func() query.PageBuilder[model.Weather], pageSize int) []model.WeatherKey {
	t.Helper()

	var keys []model.WeatherKey
	cursor := ""

	for {
		page, err := paginate().First(pageSize).After(cursor).Get(ctx)
		assert.NilError(t, err)

		keys = append(keys, weatherKeys(page.Items)...)

		if !page.PageInfo.HasNextPage {
			return keys
		}

		cursor = page.NextCursor()
	}
}

func TestPaginateArrayID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	readings := seedWeather(ctx, t, client)
	allKeys := readingKeys(readings)

	t.Run("forward by id walks every record in id order", func(t *testing.T) {
		keys := walkWeatherPages(ctx, t, client.WeatherRepo().Query().Paginate, 3)
		assert.DeepEqual(t, allKeys, keys)
	})

	t.Run("backward by id", func(t *testing.T) {
		last, err := client.WeatherRepo().Query().Paginate().Last(3).Get(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, allKeys[5:], weatherKeys(last.Items))
		assert.Equal(t, true, last.PageInfo.HasPreviousPage)

		previous, err := client.WeatherRepo().Query().Paginate().Last(3).Before(last.PrevCursor()).Get(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, allKeys[2:5], weatherKeys(previous.Items))
	})

	t.Run("custom order uses the complex id as tiebreaker", func(t *testing.T) {
		expected := slices.Clone(readings)
		slices.SortStableFunc(expected, func(left, right weatherReading) int {
			return cmp.Compare(right.temperature, left.temperature)
		})

		q := client.WeatherRepo().Query().Order(by.Weather.Temperature.Desc())

		keys := walkWeatherPages(ctx, t, q.Paginate, 2)
		assert.DeepEqual(t, readingKeys(expected), keys)
	})

	t.Run("range bounds are kept while paging", func(t *testing.T) {
		q := client.WeatherRepo().Query().Range(
			som.From(allKeys[2]),
			som.ToInclusive(allKeys[6]),
		)

		keys := walkWeatherPages(ctx, t, q.Paginate, 2)
		assert.DeepEqual(t, allKeys[2:7], keys)
	})

	t.Run("forged cursor cannot escape the range lower bound", func(t *testing.T) {
		unbounded, err := client.WeatherRepo().Query().Paginate().First(1).Get(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, allKeys[:1], weatherKeys(unbounded.Items))

		bounded, err := client.WeatherRepo().Query().Range(
			som.From(allKeys[2]),
			som.ToEnd(),
		).Paginate().First(10).After(unbounded.NextCursor()).Get(ctx)
		assert.NilError(t, err)
		assert.DeepEqual(t, allKeys[2:], weatherKeys(bounded.Items))
	})

	t.Run("total count ignores the cursor", func(t *testing.T) {
		first, err := client.WeatherRepo().Query().Paginate().First(3).Get(ctx)
		assert.NilError(t, err)

		next, err := client.WeatherRepo().Query().Paginate().First(3).After(first.NextCursor()).WithTotalCount().Get(ctx)
		assert.NilError(t, err)
		assert.Equal(t, len(allKeys), next.TotalCount)
	})
}

func TestPaginateObjectID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for _, key := range []model.PersonKey{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Diana", Age: 40},
		{Name: "Eve", Age: 45},
	} {
		person := &model.PersonObj{Node: som.NewNode[model.PersonKey](key)}
		assert.NilError(t, client.PersonObjRepo().CreateWithID(ctx, person))
	}

	expected, err := client.PersonObjRepo().Query().All(ctx)
	assert.NilError(t, err)

	var keys []model.PersonKey
	cursor := ""

	for {
		page, err := client.PersonObjRepo().Query().Paginate().First(2).After(cursor).Get(ctx)
		assert.NilError(t, err)

		for _, person := range page.Items {
			keys = append(keys, person.ID())
		}

		if !page.PageInfo.HasNextPage {
			break
		}

		cursor = page.NextCursor()
	}

	assert.Equal(t, len(expected), len(keys))

	sortedKeys := slices.Clone(keys)
	slices.SortFunc(sortedKeys, func(left, right model.PersonKey) int {
		return cmp.Or(cmp.Compare(left.Name, right.Name), cmp.Compare(left.Age, right.Age))
	})
	assert.Assert(t, len(slices.Compact(sortedKeys)) == len(keys), "pages must not repeat records")
}

func TestIterateArrayID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	readings := seedWeather(ctx, t, client)

	var keys []model.WeatherKey

	for weather, err := range client.WeatherRepo().Query().Iterate(ctx, 3) {
		assert.NilError(t, err)

		keys = append(keys, weather.ID())
	}

	assert.DeepEqual(t, readingKeys(readings), keys)
}
