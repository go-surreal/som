package basic

import (
	"context"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"som.test/gen/som"
	"som.test/gen/som/by"
	"som.test/gen/som/filter"
	"som.test/gen/som/query"
	"som.test/gen/som/with"
	"som.test/model"
)

func TestFragmentAll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := &model.AllTypes{
		FieldString:  "fragment",
		FieldInt:     42,
		FieldBool:    true,
		FieldTime:    time.Date(2026, time.March, 1, 10, 0, 0, 0, time.UTC),
		FieldEnum:    model.RoleAdmin,
		FieldFloat32: 1.5,
		FieldMonth:   time.January,
	}

	err := client.AllTypesRepo().Create(ctx, record)
	assert.NilError(t, err)

	frags, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.ID.Equal(string(record.ID()))).
		AllAs[model.AllTypesCore](ctx)
	assert.NilError(t, err)
	assert.Assert(t, is.Len(frags, 1))

	frag := frags[0]

	assert.Equal(t, frag.FieldString, "fragment")
	assert.Equal(t, frag.FieldInt, 42)
	assert.Equal(t, frag.FieldBool, true)
	assert.Equal(t, frag.FieldEnum, model.RoleAdmin)
	assert.Assert(t, frag.FieldTime.Equal(record.FieldTime))

	assert.Assert(t, frag.IsPartial(), "a fragment is always partial")
	assert.Assert(t, frag.Marker().Has(som.MarkerLoaded|som.MarkerPartial))
	assert.Equal(t, frag.RecordID(), "all_types:"+string(record.ID()))
}

func TestFragmentNarrowsSelect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	assert.Equal(t,
		client.AllTypesRepo().Query().Describe(),
		"SELECT * FROM all_types")

	assert.Equal(t,
		client.AllTypesRepo().Query().DescribeAs[model.AllTypesCore](),
		"SELECT id, field_string, field_int, field_bool, field_time, field_enum FROM all_types")

	assert.Equal(t,
		client.AllTypesRepo().Query().Order(by.AllTypes.FieldRenamed.Asc()).DescribeAs[model.AllTypesCore](),
		"SELECT id, field_string, field_int, field_bool, field_time, field_enum, custom_name FROM all_types ORDER BY custom_name ASC",
		"a sort field the fragment does not hold is added to the projection")

	record := &model.AllTypes{
		FieldString:     "narrow",
		FieldHookStatus: "not projected",
		FieldMonth:      time.January,
	}

	err := client.AllTypesRepo().Create(ctx, record)
	assert.NilError(t, err)

	frag, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldString.Equal("narrow")).
		FirstAs[model.AllTypesCore](ctx)
	assert.NilError(t, err)

	assert.Equal(t, frag.FieldString, "narrow")

	full, exists, err := client.AllTypesRepo().Expand(ctx, frag)
	assert.NilError(t, err)
	assert.Assert(t, exists)

	assert.Equal(t, full.ID(), record.ID())
	assert.Equal(t, full.FieldHookStatus, "[created]not projected",
		"the expanded record holds the fields the fragment does not project")
	assert.Assert(t, !full.IsPartial())
}

func TestFragmentRelations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	linked := &model.SpecialTypes{Name: "linked"}
	err := client.SpecialTypesRepo().Create(ctx, linked)
	assert.NilError(t, err)

	record := &model.AllTypes{
		FieldRenamed:     "renamed value",
		FieldCredentials: model.Credentials{Username: "someone"},
		FieldNodePtr:     linked,
		FieldNodeSlice:   []model.SpecialTypes{*linked},
		FieldMonth:       time.January,
	}

	err = client.AllTypesRepo().Create(ctx, record)
	assert.NilError(t, err)

	frag, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.ID.Equal(string(record.ID()))).
		FirstAs[model.AllTypesRelations](ctx)
	assert.NilError(t, err)

	assert.Equal(t, frag.FieldRenamed, "renamed value", "a renamed field is projected by its database name")
	assert.Equal(t, frag.FieldCredentials.Username, "someone")
	assert.Assert(t, frag.FieldNodePtr != nil)
	assert.Equal(t, frag.FieldNodePtr.ID(), linked.ID(), "an unfetched link keeps its record id")
	assert.Assert(t, frag.FieldNodePtr.IsPartial())
	assert.Assert(t, is.Len(frag.FieldNodeSlice, 1))
	assert.Equal(t, frag.FieldNodeSlice[0].ID(), linked.ID())

	// Fetch stays typed on the parent model and resolves the projected link.
	fetched, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.ID.Equal(string(record.ID()))).
		Fetch(with.AllTypes.FieldNodePtr()).
		FirstAs[model.AllTypesRelations](ctx)
	assert.NilError(t, err)

	assert.Assert(t, fetched.FieldNodePtr != nil)
	assert.Equal(t, fetched.FieldNodePtr.Name, "linked")
	assert.Assert(t, !fetched.FieldNodePtr.IsPartial(), "a fetched link is fully loaded")
}

func TestFragmentSortAndIterate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for _, name := range []string{"charlie", "alpha", "bravo", "delta", "echo"} {
		err := client.SpecialTypesRepo().Create(ctx, &model.SpecialTypes{Name: name})
		assert.NilError(t, err)
	}

	frags, err := client.SpecialTypesRepo().Query().
		Order(by.SpecialTypes.Name.Asc()).
		AllAs[model.SpecialTypesName](ctx)
	assert.NilError(t, err)
	assert.Assert(t, is.Len(frags, 5))
	assert.Equal(t, frags[0].Name, "alpha")
	assert.Equal(t, frags[4].Name, "echo")

	// Iteration pages through the fragments via a keyset cursor built from the
	// projected fields.
	var names []string
	for frag, err := range client.SpecialTypesRepo().Query().
		Order(by.SpecialTypes.Name.Asc()).
		IterateAs[model.SpecialTypesName](ctx, 2) {
		assert.NilError(t, err)
		names = append(names, frag.Name)
	}

	assert.DeepEqual(t, names, []string{"alpha", "bravo", "charlie", "delta", "echo"})
}

// TestFragmentSortOnUnprojectedField covers a sort on a field the fragment does
// not hold: it is added to the projection so that ORDER BY and the keyset
// cursor still work.
func TestFragmentSortOnUnprojectedField(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i, name := range []string{"first", "second", "third"} {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: name,
			FieldInt:    3 - i,
			FieldMonth:  time.January,
		})
		assert.NilError(t, err)
	}

	var names []string
	for frag, err := range client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldRenamed.Asc(), by.AllTypes.FieldInt.Asc()).
		IterateAs[model.AllTypesRelations](ctx, 2) {
		assert.NilError(t, err)
		names = append(names, frag.FieldRenamed)
	}

	assert.Assert(t, is.Len(names, 3))
}

func TestFragmentExpandComplexID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	key := model.WeatherKey{City: "berlin", Date: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)}
	record := &model.Weather{Node: som.NewNode[model.WeatherKey](key), Temperature: 12.5}

	err := client.WeatherRepo().CreateWithID(ctx, record)
	assert.NilError(t, err)

	frag, err := client.WeatherRepo().Query().
		FirstAs[model.WeatherReading](ctx)
	assert.NilError(t, err)
	assert.Equal(t, frag.Temperature, 12.5)

	full, exists, err := client.WeatherRepo().Expand(ctx, frag)
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, full.ID().City, "berlin")
	assert.Assert(t, full.ID().Date.Equal(key.Date))
}

func TestFragmentLive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	liveChan, err := client.SpecialTypesRepo().Query().LiveAs[model.SpecialTypesName](ctx)
	assert.NilError(t, err)

	record := &model.SpecialTypes{Name: "live location"}
	err = client.SpecialTypesRepo().Create(ctx, record)
	assert.NilError(t, err)

	select {
	case res, more := <-liveChan:
		assert.Assert(t, more, "liveChan closed unexpectedly")

		created, ok := res.(query.LiveCreate[*model.SpecialTypesName])
		assert.Assert(t, ok, "expected LiveCreate event, got %T", res)

		frag, err := created.Get()
		assert.NilError(t, err)
		assert.Equal(t, frag.Name, "live location")
		assert.Assert(t, frag.IsPartial())

	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for CREATE event")
	}
}

func TestFragmentIsReadOnly(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := &model.SpecialTypes{Name: "read only"}
	err := client.SpecialTypesRepo().Create(ctx, record)
	assert.NilError(t, err)

	frag, err := client.SpecialTypesRepo().Query().FirstAs[model.SpecialTypesName](ctx)
	assert.NilError(t, err)

	// A fragment cannot be written back: it is not the model type the write
	// methods accept, and expanding it yields a full model that can be.
	full, exists, err := client.SpecialTypesRepo().Expand(ctx, frag)
	assert.NilError(t, err)
	assert.Assert(t, exists)

	err = client.SpecialTypesRepo().Update(ctx, full)
	assert.NilError(t, err)
}
