package basic

import (
	"context"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestComplexIDObjectKey(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	key := model.PersonKey{Name: "Alice", Age: 30}
	person := &model.PersonObj{
		CustomNode: som.NewCustomNode[model.PersonKey](key),
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
		CustomNode:  som.NewCustomNode[model.WeatherKey](key),
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
}

func TestComplexIDMultipleRecords(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	key1 := model.PersonKey{Name: "Alice", Age: 30}
	key2 := model.PersonKey{Name: "Bob", Age: 25}

	person1 := &model.PersonObj{
		CustomNode: som.NewCustomNode[model.PersonKey](key1),
		Email:      "alice@example.com",
	}
	person2 := &model.PersonObj{
		CustomNode: som.NewCustomNode[model.PersonKey](key2),
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
		CustomNode:  som.NewCustomNode[model.WeatherKey](weatherKey),
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
		CustomNode: som.NewCustomNode[model.TeamMemberKey](tmKey),
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
