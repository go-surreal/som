package basic

import (
	"context"
	"testing"
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
	}

	userIn := userNew

	err := client.AllTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.ID.Equal(userIn.ID()),
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}, som.ID{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus"),
	)
}

func TestRefresh(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	allTypes := &model.AllTypes{
		FieldString: "some value",
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

func FuzzWithDatabase(f *testing.F) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, f)
	defer cleanup()

	f.Add("Some User")

	f.Fuzz(func(t *testing.T, str string) {
		userIn := &model.AllTypes{
			FieldString: str,
		}

		err := client.AllTypesRepo().Create(ctx, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == nil {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllTypesRepo().Query().
			Where(
				filter.AllTypes.ID.Equal(userIn.ID()),
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
		}

		err := client.AllTypesRepo().CreateWithID(ctx, id, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == nil {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, ok, err := client.AllTypesRepo().Read(ctx, userIn.ID())

		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			t.Fatalf("user with ID '%s' not found", userIn.ID())
		}

		assert.Equal(t, userIn.ID().String(), userOut.ID().String())
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
		}

		err := client.AllTypesRepo().Create(ctx, userIn)
		if err != nil {
			b.Fatal(err)
		}

		if userIn.ID() == nil {
			b.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllTypesRepo().Query().
			Where(
				filter.AllTypes.ID.Equal(userIn.ID()),
			).
			First(ctx)

		if err != nil {
			b.Fatal(err)
		}

		assert.Equal(b, userIn.FieldString, userOut.FieldString)
	}
}
