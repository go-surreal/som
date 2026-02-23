package basic

import (
	"testing"
	"time"

	som "github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestRangeQuery(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	t.Run("ArrayID full range", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.From(model.WeatherKey{City: "London", Date: start}),
			som.To(model.WeatherKey{City: "London", Date: end}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM weather:[$A, $B]..[$C, $D]",
			desc,
		)
	})

	t.Run("ArrayID open-ended to end", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.From(model.WeatherKey{City: "London", Date: start}),
			som.ToEnd(),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM weather:[$A, $B]..",
			desc,
		)
	})

	t.Run("ArrayID open-ended from start", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.FromStart(),
			som.To(model.WeatherKey{City: "London", Date: end}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM weather..[$A, $B]",
			desc,
		)
	})

	t.Run("ArrayID exclusive from inclusive to", func(t *testing.T) {
		desc := query.NewWeather(nil).Range(
			som.FromExclusive(model.WeatherKey{City: "London", Date: start}),
			som.ToInclusive(model.WeatherKey{City: "London", Date: end}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM weather:[$A, $B]>..=[$C, $D]",
			desc,
		)
	})

	t.Run("ObjectID full range", func(t *testing.T) {
		desc := query.NewPersonObj(nil).Range(
			som.From(model.PersonKey{Name: "Alice", Age: 20}),
			som.To(model.PersonKey{Name: "Bob", Age: 30}),
		).Describe()

		assert.Equal(t,
			"SELECT * FROM person_obj:{name: $A, age: $B}..{name: $C, age: $D}",
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
			"SELECT * FROM weather:[$A, $B]..[$C, $D] WHERE (temperature > $E)",
			desc,
		)
	})
}
