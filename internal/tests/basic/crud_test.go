package basic

import (
	"context"
	"regexp"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestWithDatabase(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "Some User"
	uid := uuid.New()

	userNew := model.AllTypes{
		FieldString:    str,
		FieldUUID:      uid,
		FieldByte:      []byte("x")[0],
		FieldByteSlice: []byte("some value"),
		FieldMonth:     time.March,
		FieldWeekday:   time.Wednesday,
	}

	userIn := userNew

	err := client.AllTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldMonth.Equal(time.March),
			filter.AllTypes.ID.Equal(string(userIn.ID())),
			filter.AllTypes.FieldString.Equal(str),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, str, userOut.FieldString)
	assert.Equal(t, uid, userOut.FieldUUID)
	assert.Equal(t, userNew.FieldByte, userOut.FieldByte)
	assert.DeepEqual(t, userNew.FieldByteSlice, userOut.FieldByteSlice)

	assert.DeepEqual(t,
		userNew, *userOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}, regexp.Regexp{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus"),
	)
}

func TestMonthWeekdayPointers(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	march := time.March
	wednesday := time.Wednesday

	// Create with non-nil pointer values.
	withValues := &model.AllTypes{
		FieldMonth:      time.January,
		FieldMonthPtr:   &march,
		FieldWeekday:    time.Monday,
		FieldWeekdayPtr: &wednesday,
	}

	err := client.AllTypesRepo().Create(ctx, withValues)
	if err != nil {
		t.Fatal(err)
	}

	readWithValues, ok, err := client.AllTypesRepo().Read(ctx, string(withValues.ID()))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, ok)
	assert.Assert(t, readWithValues.FieldMonthPtr != nil)
	assert.Equal(t, time.March, *readWithValues.FieldMonthPtr)
	assert.Assert(t, readWithValues.FieldWeekdayPtr != nil)
	assert.Equal(t, time.Wednesday, *readWithValues.FieldWeekdayPtr)

	// Create with nil pointer values.
	withNil := &model.AllTypes{
		FieldMonth: time.January,
	}

	err = client.AllTypesRepo().Create(ctx, withNil)
	if err != nil {
		t.Fatal(err)
	}

	readWithNil, ok, err := client.AllTypesRepo().Read(ctx, string(withNil.ID()))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, ok)
	assert.Assert(t, readWithNil.FieldMonthPtr == nil)
	assert.Assert(t, readWithNil.FieldWeekdayPtr == nil)

	// Filter: non-nil pointer values.
	results, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldMonthPtr.Equal(time.March)).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, withValues.ID(), results[0].ID())

	// Filter: nil pointer values.
	results, err = client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldMonthPtr.Nil(true)).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, withNil.ID(), results[0].ID())

	// Filter: weekday pointer.
	results, err = client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldWeekdayPtr.Equal(time.Wednesday)).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, withValues.ID(), results[0].ID())

	results, err = client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldWeekdayPtr.Nil(true)).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, withNil.ID(), results[0].ID())
}

func TestRefresh(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	allTypes := &model.AllTypes{
		FieldString: "some value",
		FieldMonth:  time.January,
	}

	err := client.AllTypesRepo().Create(ctx, allTypes)
	if err != nil {
		t.Fatal(err)
	}

	allTypes.FieldString = "some other value"

	err = client.AllTypesRepo().Refresh(ctx, allTypes)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "some value", allTypes.FieldString)
}

func TestInsert(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	nodes := []*model.AllTypes{
		{FieldString: "first", FieldMonth: time.January},
		{FieldString: "second", FieldMonth: time.January},
		{FieldString: "third", FieldMonth: time.January},
	}

	err := client.AllTypesRepo().Insert(ctx, nodes)
	if err != nil {
		t.Fatal(err)
	}

	for i, n := range nodes {
		if n.ID() == "" {
			t.Fatalf("node %d: ID must not be empty after insert", i)
		}
	}

	assert.Equal(t, "first", nodes[0].FieldString)
	assert.Equal(t, "second", nodes[1].FieldString)
	assert.Equal(t, "third", nodes[2].FieldString)

	// Verify records exist in the database.
	for _, n := range nodes {
		out, ok, err := client.AllTypesRepo().Read(ctx, string(n.ID()))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("record %s not found after insert", n.ID())
		}
		assert.Equal(t, n.FieldString, out.FieldString)
	}
}

func TestInsertEmpty(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Insert(ctx, nil)
	assert.NilError(t, err)

	err = client.AllTypesRepo().Insert(ctx, []*model.AllTypes{})
	assert.NilError(t, err)
}

func TestInsertValidation(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("nil element", func(t *testing.T) {
		err := client.AllTypesRepo().Insert(ctx, []*model.AllTypes{nil})
		assert.ErrorContains(t, err, "nil node")
	})

	t.Run("node with existing id", func(t *testing.T) {
		n := &model.AllTypes{FieldString: "test", FieldMonth: time.January}
		err := client.AllTypesRepo().Create(ctx, n)
		if err != nil {
			t.Fatal(err)
		}
		err = client.AllTypesRepo().Insert(ctx, []*model.AllTypes{n})
		assert.ErrorContains(t, err, "already has an id")
	})
}

func FuzzWithDatabase(f *testing.F) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, f)
	defer cleanup()

	f.Add("Some User")

	f.Fuzz(func(t *testing.T, str string) {
		userIn := &model.AllTypes{
			FieldString: str,
			FieldMonth:  time.January,
		}

		err := client.AllTypesRepo().Create(ctx, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == "" {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllTypesRepo().Query().
			Where(
				filter.AllTypes.ID.Equal(string(userIn.ID())),
			).
			First(ctx)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, userIn.FieldString, userOut.FieldString)
	})
}

func FuzzCustomModelIDs(f *testing.F) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, f)
	defer cleanup()

	f.Add("v9uitj942tv2403tnv")
	f.Add("vb92thj29v4tjn20d3")
	f.Add("ij024itvnjc20394it")
	f.Add(" 0")
	f.Add("\"0")
	f.Add("🙂")
	f.Add("✅")
	f.Add("👋😉")

	f.Fuzz(func(t *testing.T, id string) {
		if !utf8.ValidString(id) {
			t.Skip("id is not a valid utf8 string")
		}

		userIn := &model.AllTypes{
			FieldString: "1",
			FieldMonth:  time.January,
		}

		err := client.AllTypesRepo().CreateWithID(ctx, id, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == "" {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, ok, err := client.AllTypesRepo().Read(ctx, string(userIn.ID()))

		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			t.Fatalf("user with ID '%s' not found", userIn.ID())
		}

		assert.Equal(t, userIn.ID(), userOut.ID())
		assert.Equal(t, "1", userOut.FieldString)

		userOut.FieldString = "2"

		err = client.AllTypesRepo().Update(ctx, userOut)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "2", userOut.FieldString)

		err = client.AllTypesRepo().Delete(ctx, userOut)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func BenchmarkWithDatabase(b *testing.B) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, b)
	defer cleanup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		userIn := &model.AllTypes{
			FieldString: "Some User",
			FieldMonth:  time.January,
		}

		err := client.AllTypesRepo().Create(ctx, userIn)
		if err != nil {
			b.Fatal(err)
		}

		if userIn.ID() == "" {
			b.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllTypesRepo().Query().
			Where(
				filter.AllTypes.ID.Equal(string(userIn.ID())),
			).
			First(ctx)

		if err != nil {
			b.Fatal(err)
		}

		assert.Equal(b, userIn.FieldString, userOut.FieldString)
	}
}
