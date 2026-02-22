package basic

import (
	"context"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestComplexIDObjectKey(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	key := model.PersonKey{Name: "Alice", Age: 30}
	person := &model.PersonObj{
		Node: som.NewNode[model.PersonKey](key),
		Email:      "alice@example.com",
	}

	// Create
	err := client.PersonObjRepo().CreateWithID(ctx, person)
	assert.NilError(t, err)

	// Read
	read, ok, err := client.PersonObjRepo().Read(ctx, key)
	assert.NilError(t, err)
	assert.Assert(t, ok, "expected record to exist")
	assert.Equal(t, read.Email, "alice@example.com")
	assert.Equal(t, read.ID().Name, "Alice")
	assert.Equal(t, read.ID().Age, 30)

	// Update
	person.Email = "alice-updated@example.com"
	err = client.PersonObjRepo().Update(ctx, person)
	assert.NilError(t, err)

	read, ok, err = client.PersonObjRepo().Read(ctx, key)
	assert.NilError(t, err)
	assert.Assert(t, ok)
	assert.Equal(t, read.Email, "alice-updated@example.com")

	// Refresh
	person.Email = "should-be-overwritten"
	err = client.PersonObjRepo().Refresh(ctx, person)
	assert.NilError(t, err)
	assert.Equal(t, person.Email, "alice-updated@example.com")

	// Delete
	err = client.PersonObjRepo().Delete(ctx, person)
	assert.NilError(t, err)

	_, ok, err = client.PersonObjRepo().Read(ctx, key)
	assert.NilError(t, err)
	assert.Assert(t, !ok, "expected record to be deleted")
}

func TestComplexIDArrayKey(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	fixedDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	key := model.WeatherKey{City: "Berlin", Date: fixedDate}
	w := &model.Weather{
		Node:  som.NewNode[model.WeatherKey](key),
		Temperature: 22.5,
	}

	// Create
	err := client.WeatherRepo().CreateWithID(ctx, w)
	assert.NilError(t, err)

	// Read
	read, ok, err := client.WeatherRepo().Read(ctx, key)
	assert.NilError(t, err)
	assert.Assert(t, ok, "expected record to exist")
	assert.Equal(t, read.Temperature, 22.5)
	assert.Equal(t, read.ID().City, "Berlin")
	assert.Assert(t, read.ID().Date.Equal(fixedDate), "expected date to match")

	// Update
	w.Temperature = 25.0
	err = client.WeatherRepo().Update(ctx, w)
	assert.NilError(t, err)

	read, ok, err = client.WeatherRepo().Read(ctx, key)
	assert.NilError(t, err)
	assert.Assert(t, ok)
	assert.Equal(t, read.Temperature, 25.0)

	// Refresh
	w.Temperature = 999.0
	err = client.WeatherRepo().Refresh(ctx, w)
	assert.NilError(t, err)
	assert.Equal(t, w.Temperature, 25.0)

	// Delete
	err = client.WeatherRepo().Delete(ctx, w)
	assert.NilError(t, err)

	_, ok, err = client.WeatherRepo().Read(ctx, key)
	assert.NilError(t, err)
	assert.Assert(t, !ok, "expected record to be deleted")
}

func TestComplexIDZeroKeyErrors(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("create_zero_person_key", func(t *testing.T) {
		person := &model.PersonObj{}
		err := client.PersonObjRepo().CreateWithID(ctx, person)
		assert.ErrorContains(t, err, "non-zero ID")
	})

	t.Run("update_zero_person_key", func(t *testing.T) {
		person := &model.PersonObj{}
		err := client.PersonObjRepo().Update(ctx, person)
		assert.ErrorContains(t, err, "without existing record ID")
	})

	t.Run("delete_zero_person_key", func(t *testing.T) {
		person := &model.PersonObj{}
		err := client.PersonObjRepo().Delete(ctx, person)
		assert.ErrorContains(t, err, "without existing record ID")
	})

	t.Run("create_zero_weather_key", func(t *testing.T) {
		w := &model.Weather{}
		err := client.WeatherRepo().CreateWithID(ctx, w)
		assert.ErrorContains(t, err, "non-zero ID")
	})

	t.Run("update_zero_weather_key", func(t *testing.T) {
		w := &model.Weather{}
		err := client.WeatherRepo().Update(ctx, w)
		assert.ErrorContains(t, err, "without existing record ID")
	})

	t.Run("delete_zero_weather_key", func(t *testing.T) {
		w := &model.Weather{}
		err := client.WeatherRepo().Delete(ctx, w)
		assert.ErrorContains(t, err, "without existing record ID")
	})

	t.Run("create_zero_team_member_key", func(t *testing.T) {
		tm := &model.TeamMember{}
		err := client.TeamMemberRepo().CreateWithID(ctx, tm)
		assert.ErrorContains(t, err, "Member.ID must not be empty")
	})

	t.Run("update_zero_team_member_key", func(t *testing.T) {
		tm := &model.TeamMember{}
		err := client.TeamMemberRepo().Update(ctx, tm)
		assert.ErrorContains(t, err, "Member.ID must not be empty")
	})

	t.Run("delete_zero_team_member_key", func(t *testing.T) {
		tm := &model.TeamMember{}
		err := client.TeamMemberRepo().Delete(ctx, tm)
		assert.ErrorContains(t, err, "Member.ID must not be empty")
	})
}

func TestComplexIDMultipleRecords(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	key1 := model.PersonKey{Name: "Alice", Age: 30}
	key2 := model.PersonKey{Name: "Bob", Age: 25}

	person1 := &model.PersonObj{
		Node: som.NewNode[model.PersonKey](key1),
		Email:      "alice@example.com",
	}
	person2 := &model.PersonObj{
		Node: som.NewNode[model.PersonKey](key2),
		Email:      "bob@example.com",
	}

	err := client.PersonObjRepo().CreateWithID(ctx, person1)
	assert.NilError(t, err)

	err = client.PersonObjRepo().CreateWithID(ctx, person2)
	assert.NilError(t, err)

	read1, ok, err := client.PersonObjRepo().Read(ctx, key1)
	assert.NilError(t, err)
	assert.Assert(t, ok)
	assert.Equal(t, read1.Email, "alice@example.com")

	read2, ok, err := client.PersonObjRepo().Read(ctx, key2)
	assert.NilError(t, err)
	assert.Assert(t, ok)
	assert.Equal(t, read2.Email, "bob@example.com")
}

func TestComplexIDNodeRef(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Create referenced nodes first.
	member := &model.AllTypes{FieldString: "ref-member"}
	err := client.AllTypesRepo().Create(ctx, member)
	assert.NilError(t, err)
	assert.Assert(t, member.ID() != "")

	fixedDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	weatherKey := model.WeatherKey{City: "Berlin", Date: fixedDate}
	weather := &model.Weather{
		Node:  som.NewNode[model.WeatherKey](weatherKey),
		Temperature: 22.5,
	}
	err = client.WeatherRepo().CreateWithID(ctx, weather)
	assert.NilError(t, err)

	// Create TeamMember with node references in key.
	tmKey := model.TeamMemberKey{
		Member:   *member,
		Forecast: *weather,
	}
	tm := &model.TeamMember{
		Node: som.NewNode[model.TeamMemberKey](tmKey),
		Role:       "engineer",
	}
	err = client.TeamMemberRepo().CreateWithID(ctx, tm)
	assert.NilError(t, err)

	// Read back and verify.
	read, ok, err := client.TeamMemberRepo().Read(ctx, tmKey)
	assert.NilError(t, err)
	assert.Assert(t, ok, "expected record to exist")
	assert.Equal(t, read.Role, "engineer")
	assert.Equal(t, read.ID().Member.ID(), member.ID())
	assert.Equal(t, read.ID().Forecast.ID().City, "Berlin")
	assert.Assert(t, read.ID().Forecast.ID().Date.Equal(fixedDate))

	// Update.
	tm.Role = "senior-engineer"
	err = client.TeamMemberRepo().Update(ctx, tm)
	assert.NilError(t, err)

	read, ok, err = client.TeamMemberRepo().Read(ctx, tmKey)
	assert.NilError(t, err)
	assert.Assert(t, ok)
	assert.Equal(t, read.Role, "senior-engineer")

	// Refresh.
	tm.Role = "should-be-overwritten"
	err = client.TeamMemberRepo().Refresh(ctx, tm)
	assert.NilError(t, err)
	assert.Equal(t, tm.Role, "senior-engineer")

	// Delete.
	err = client.TeamMemberRepo().Delete(ctx, tm)
	assert.NilError(t, err)

	_, ok, err = client.TeamMemberRepo().Read(ctx, tmKey)
	assert.NilError(t, err)
	assert.Assert(t, !ok, "expected record to be deleted")
}

func TestComplexIDQueryAll(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	keys := []model.PersonKey{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	emails := []string{"alice@example.com", "bob@example.com", "charlie@example.com"}

	for i, key := range keys {
		p := &model.PersonObj{
			Node: som.NewNode[model.PersonKey](key),
			Email:      emails[i],
		}
		err := client.PersonObjRepo().CreateWithID(ctx, p)
		assert.NilError(t, err)
	}

	all, err := client.PersonObjRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(all), 3)

	count, err := client.PersonObjRepo().Query().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, count, 3)
}

func TestComplexIDQueryFilter(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	keys := []model.PersonKey{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	emails := []string{"alice@example.com", "bob@example.com", "charlie@example.com"}

	for i, key := range keys {
		p := &model.PersonObj{
			Node: som.NewNode[model.PersonKey](key),
			Email:      emails[i],
		}
		err := client.PersonObjRepo().CreateWithID(ctx, p)
		assert.NilError(t, err)
	}

	results, err := client.PersonObjRepo().Query().
		Where(filter.PersonObj.Email.Equal("bob@example.com")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Email, "bob@example.com")

	results, err = client.PersonObjRepo().Query().
		Where(filter.PersonObj.Email.Equal("nonexistent@example.com")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(results), 0)
}

func TestComplexIDQuerySort(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	temps := []float64{22.5, 10.0, 35.0}
	cities := []string{"Berlin", "London", "Tokyo"}

	for i, city := range cities {
		fixedDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
		key := model.WeatherKey{City: city, Date: fixedDate}
		w := &model.Weather{
			Node:  som.NewNode[model.WeatherKey](key),
			Temperature: temps[i],
		}
		err := client.WeatherRepo().CreateWithID(ctx, w)
		assert.NilError(t, err)
	}

	asc, err := client.WeatherRepo().Query().
		Order(by.Weather.Temperature.Asc()).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(asc), 3)
	assert.Equal(t, asc[0].Temperature, 10.0)
	assert.Equal(t, asc[1].Temperature, 22.5)
	assert.Equal(t, asc[2].Temperature, 35.0)

	desc, err := client.WeatherRepo().Query().
		Order(by.Weather.Temperature.Desc()).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(desc), 3)
	assert.Equal(t, desc[0].Temperature, 35.0)
	assert.Equal(t, desc[1].Temperature, 22.5)
	assert.Equal(t, desc[2].Temperature, 10.0)
}

func TestComplexIDQueryLimitOffset(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	keys := []model.PersonKey{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	emails := []string{"alice@example.com", "bob@example.com", "charlie@example.com"}

	for i, key := range keys {
		p := &model.PersonObj{
			Node: som.NewNode[model.PersonKey](key),
			Email:      emails[i],
		}
		err := client.PersonObjRepo().CreateWithID(ctx, p)
		assert.NilError(t, err)
	}

	limited, err := client.PersonObjRepo().Query().
		Order(by.PersonObj.Email.Asc()).
		Limit(2).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(limited), 2)
	assert.Equal(t, limited[0].Email, "alice@example.com")
	assert.Equal(t, limited[1].Email, "bob@example.com")

	offset, err := client.PersonObjRepo().Query().
		Order(by.PersonObj.Email.Asc()).
		Offset(1).
		Limit(1).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(offset), 1)
	assert.Equal(t, offset[0].Email, "bob@example.com")
}

func TestComplexIDQueryFirst(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	fixedDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	for _, city := range []string{"Berlin", "London"} {
		key := model.WeatherKey{City: city, Date: fixedDate}
		temp := 22.5
		if city == "London" {
			temp = 10.0
		}
		w := &model.Weather{
			Node:  som.NewNode[model.WeatherKey](key),
			Temperature: temp,
		}
		err := client.WeatherRepo().CreateWithID(ctx, w)
		assert.NilError(t, err)
	}

	first, err := client.WeatherRepo().Query().
		Order(by.Weather.Temperature.Asc()).
		First(ctx)
	assert.NilError(t, err)
	assert.Equal(t, first.Temperature, 10.0)
	assert.Equal(t, first.ID().City, "London")
}

func TestComplexIDFilterByObjectIDField(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	keys := []model.PersonKey{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	emails := []string{"alice@example.com", "bob@example.com", "charlie@example.com"}

	for i, key := range keys {
		p := &model.PersonObj{
			Node: som.NewNode[model.PersonKey](key),
			Email:      emails[i],
		}
		err := client.PersonObjRepo().CreateWithID(ctx, p)
		assert.NilError(t, err)
	}

	// Filter by ID sub-field: Name
	results, err := client.PersonObjRepo().Query().
		Where(filter.PersonObj.ID().Name.Equal("Alice")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Email, "alice@example.com")

	// Filter by ID sub-field: Age > 26 → Alice(30) and Charlie(35)
	results, err = client.PersonObjRepo().Query().
		Where(filter.PersonObj.ID().Age.GreaterThan(26)).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(results), 2)
}

func TestComplexIDFilterByArrayIDField(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	dates := []time.Time{
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
	}
	cities := []string{"Berlin", "London", "Tokyo"}
	temps := []float64{5.0, 20.0, 30.0}

	for i, city := range cities {
		key := model.WeatherKey{City: city, Date: dates[i]}
		w := &model.Weather{
			Node:  som.NewNode[model.WeatherKey](key),
			Temperature: temps[i],
		}
		err := client.WeatherRepo().CreateWithID(ctx, w)
		assert.NilError(t, err)
	}

	// Filter by ID sub-field: City
	results, err := client.WeatherRepo().Query().
		Where(filter.Weather.ID().City.Equal("Berlin")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Temperature, 5.0)

	// Filter by ID sub-field: Date >= June 1st → London and Tokyo
	cutoff := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	results, err = client.WeatherRepo().Query().
		Where(filter.Weather.ID().Date.AfterOrEqual(cutoff)).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(results), 2)
}

func TestComplexIDQueryNodeRef(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	member1 := &model.AllTypes{FieldString: "member-1"}
	member2 := &model.AllTypes{FieldString: "member-2"}
	err := client.AllTypesRepo().Create(ctx, member1)
	assert.NilError(t, err)
	err = client.AllTypesRepo().Create(ctx, member2)
	assert.NilError(t, err)

	fixedDate := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	weather := &model.Weather{
		Node:  som.NewNode[model.WeatherKey](model.WeatherKey{City: "Berlin", Date: fixedDate}),
		Temperature: 22.5,
	}
	err = client.WeatherRepo().CreateWithID(ctx, weather)
	assert.NilError(t, err)

	tm1 := &model.TeamMember{
		Node: som.NewNode[model.TeamMemberKey](model.TeamMemberKey{
			Member:   *member1,
			Forecast: *weather,
		}),
		Role: "engineer",
	}
	tm2 := &model.TeamMember{
		Node: som.NewNode[model.TeamMemberKey](model.TeamMemberKey{
			Member:   *member2,
			Forecast: *weather,
		}),
		Role: "designer",
	}
	err = client.TeamMemberRepo().CreateWithID(ctx, tm1)
	assert.NilError(t, err)
	err = client.TeamMemberRepo().CreateWithID(ctx, tm2)
	assert.NilError(t, err)

	results, err := client.TeamMemberRepo().Query().
		Where(filter.TeamMember.Role.Equal("engineer")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Role, "engineer")

	all, err := client.TeamMemberRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(all), 2)
}
