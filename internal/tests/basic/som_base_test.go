package basic

import (
	"context"
	"math"
	"net/url"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/repo"
	"github.com/go-surreal/som/tests/basic/gen/som/where"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestQuery(t *testing.T) {
	client := &repo.ClientImpl{}

	query := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.
				MemberOf(
					where.GroupMember.CreatedAt.Before(time.Now()),
				).
				Group(
					where.Group.ID.Equal(som.MakeID("all_field_types", "some_id")),
				),

			where.AllFieldTypes.Duration.Days().LessThan(4),

			//where.AllFieldTypes.Float64.Equal_(constant.E[model.AllFieldTypes]()),
			//
			//constant.String[model.AllFieldTypes]("A").Equal_(constant.String[model.AllFieldTypes]("A")),
		)

	assert.Equal(t,
		"SELECT * OMIT login.password, login.password_ptr FROM all_field_types WHERE (->group_member[WHERE (created_at < $A)]->group[WHERE (id = $B)] "+
			"AND duration::days(duration) < $C)",
		query.Describe(),
	)

	query = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StringPtr.Base64Decode().Base64Encode().
				Equal_(where.AllFieldTypes.String.Base64Decode().Base64Encode()),
		)

	assert.Equal(t,
		"SELECT * OMIT login.password, login.password_ptr FROM all_field_types WHERE "+
			"(encoding::base64::encode(encoding::base64::decode(string_ptr)) "+
			"= encoding::base64::encode(encoding::base64::decode(string)))",
		query.Describe(),
	)
}

func TestWithDatabase(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "Some User"
	uid := uuid.New()

	userNew := model.AllFieldTypes{
		String:    str,
		UUID:      uid,
		Byte:      []byte("x")[0],
		ByteSlice: []byte("some value"),
	}

	userIn := userNew

	err := client.AllFieldTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(userIn.ID()),
			where.AllFieldTypes.String.Equal(str),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, str, userOut.String)
	assert.Equal(t, uid, userOut.UUID)
	assert.Equal(t, userNew.Byte, userOut.Byte)
	assert.DeepEqual(t, userNew.ByteSlice, userOut.ByteSlice)

	assert.DeepEqual(t,
		userNew, *userOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)
}

func TestNumerics(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "user"

	// MAX

	userMax := model.AllFieldTypes{
		String: str,
		Int:    math.MaxInt,
		Int8:   math.MaxInt8,
		Int16:  math.MaxInt16,
		Int32:  math.MaxInt32,
		Int64:  math.MaxInt64,
		//Uint:    1, //math.MaxUint,
		Uint8:  math.MaxUint8,
		Uint16: math.MaxUint16,
		Uint32: math.MaxUint32,
		//Uint64:  1, //math.MaxUint64,
		//Uintptr: 1, //math.MaxUint64,
		Float32: math.MaxFloat32,
		Float64: math.MaxFloat64,
		Rune:    math.MaxInt32,
	}

	userIn := userMax

	err := client.AllFieldTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(userIn.ID()),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t,
		userMax, *userOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)

	// MIN

	userMin := model.AllFieldTypes{
		String: str,
		Int:    math.MinInt,
		Int8:   math.MinInt8,
		Int16:  math.MinInt16,
		Int32:  math.MinInt32,
		Int64:  math.MinInt64,
		//Uint:    math.MaxUint,
		Uint8:  0,
		Uint16: 0,
		Uint32: 0,
		//Uint64:  0,
		//Uintptr: 0,
		Float32: -math.MaxFloat32,
		Float64: -math.MaxFloat64,
		Rune:    math.MinInt32,
	}

	userIn = userMin

	err = client.AllFieldTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(userIn.ID()),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t,
		userMin, *userOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)
}

func TestSlice(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// initial nil slice

	user := &model.AllFieldTypes{}

	err := client.AllFieldTypesRepo().Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.StructSlice == nil)

	user, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StructSlice.IsEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.StructSlice == nil)

	// empty slice

	user.StructSlice = []model.SomeStruct{}

	assert.Check(t, user.StructSlice != nil)

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	user, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StructSlice.Empty(true),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.StructSlice != nil)

	// non-empty slice with actual data

	str1 := "hello"
	num1 := 42
	now1 := time.Now().Truncate(time.Microsecond).UTC()
	id1 := uuid.New()

	user.StructSlice = []model.SomeStruct{{
		StringPtr: &str1,
		IntPtr:    &num1,
		TimePtr:   &now1,
		UuidPtr:   &id1,
	}}

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity: %v", err)
	}

	user, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StructSlice.NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.StructSlice) == 1)
	assert.Check(t, user.StructSlice[0].StringPtr != nil && *user.StructSlice[0].StringPtr == str1)
	assert.Check(t, user.StructSlice[0].IntPtr != nil && *user.StructSlice[0].IntPtr == num1)
	assert.Check(t, user.StructSlice[0].TimePtr != nil && user.StructSlice[0].TimePtr.Equal(now1))
	assert.Check(t, user.StructSlice[0].UuidPtr != nil && *user.StructSlice[0].UuidPtr == id1)

	// multiple elements

	str2 := "world"
	num2 := 99

	user.StructSlice = []model.SomeStruct{
		{StringPtr: &str1, IntPtr: &num1},
		{StringPtr: &str2, IntPtr: &num2},
	}

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity with multiple elements: %v", err)
	}

	user, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StructSlice.NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.StructSlice) == 2)
	assert.Check(t, *user.StructSlice[0].StringPtr == str1)
	assert.Check(t, *user.StructSlice[0].IntPtr == num1)
	assert.Check(t, *user.StructSlice[1].StringPtr == str2)
	assert.Check(t, *user.StructSlice[1].IntPtr == num2)

	// test StructPtrSlice ([]*SomeStruct)

	user.StructPtrSlice = []*model.SomeStruct{
		{StringPtr: &str1, IntPtr: &num1},
		{StringPtr: &str2, IntPtr: &num2},
	}

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity with StructPtrSlice: %v", err)
	}

	user, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StructPtrSlice.NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.StructPtrSlice) == 2)
	assert.Check(t, user.StructPtrSlice[0] != nil && *user.StructPtrSlice[0].StringPtr == str1)
	assert.Check(t, user.StructPtrSlice[1] != nil && *user.StructPtrSlice[1].StringPtr == str2)

	// test StructPtrSlicePtr (*[]*SomeStruct)

	ptrSlice := []*model.SomeStruct{
		{StringPtr: &str1, IntPtr: &num1},
	}
	user.StructPtrSlicePtr = &ptrSlice

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity with StructPtrSlicePtr: %v", err)
	}

	user, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StructPtrSlicePtr.NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.StructPtrSlicePtr != nil)
	assert.Check(t, len(*user.StructPtrSlicePtr) == 1)
	assert.Check(t, (*user.StructPtrSlicePtr)[0] != nil && *(*user.StructPtrSlicePtr)[0].StringPtr == str1)

	// test refresh with struct slice data

	str3 := "refreshed"
	user.StructSlice = []model.SomeStruct{{StringPtr: &str3}}

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity for refresh test: %v", err)
	}

	// modify local data
	modified := "modified"
	user.StructSlice[0].StringPtr = &modified

	// refresh should restore to DB value
	err = client.AllFieldTypesRepo().Refresh(ctx, user)
	if err != nil {
		t.Fatalf("could not refresh entity: %v", err)
	}

	assert.Check(t, len(user.StructSlice) == 1)
	assert.Check(t, *user.StructSlice[0].StringPtr == str3)
}

func TestTimestamps(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user := &model.AllFieldTypes{}

	err := client.AllFieldTypesRepo().Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, !user.CreatedAt().IsZero())
	assert.Check(t, !user.UpdatedAt().IsZero())
	assert.Check(t, time.Since(user.CreatedAt()) < time.Second)
	assert.Check(t, time.Since(user.UpdatedAt()) < time.Second)

	time.Sleep(time.Second)

	err = client.AllFieldTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, !user.CreatedAt().IsZero())
	assert.Check(t, !user.UpdatedAt().IsZero())
	assert.Check(t, time.Since(user.CreatedAt()) > time.Second)
	assert.Check(t, time.Since(user.UpdatedAt()) < time.Second)
}

func TestURLTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	someURL, err := url.Parse("https://surrealdb.com")
	if err != nil {
		t.Fatal(err)
	}

	newModel := &model.URLExample{
		SomeURL:      someURL,
		SomeOtherURL: *someURL,
	}

	err = client.URLExampleRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	readModel, exists, err := client.URLExampleRepo().Read(ctx, newModel.ID())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, exists)

	assert.Equal(t, someURL.String(), readModel.SomeURL.String())
	assert.Equal(t, someURL.String(), readModel.SomeOtherURL.String())

	someURL, err = url.Parse("https://github.com/surrealdb/surrealdb")
	if err != nil {
		t.Fatal(err)
	}

	readModel.SomeURL = someURL
	readModel.SomeOtherURL = *someURL

	err = client.URLExampleRepo().Update(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, someURL.String(), readModel.SomeURL.String())
	assert.Equal(t, someURL.String(), readModel.SomeOtherURL.String())

	queryModel, err := client.URLExampleRepo().Query().
		Filter(
			where.URLExample.SomeURL.Equal(*someURL),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, someURL.String(), queryModel.SomeURL.String())

	err = client.URLExampleRepo().Delete(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDuration(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ptr := time.Hour

	userNew := &model.AllFieldTypes{
		Duration:    time.Minute,
		DurationPtr: &ptr,
		DurationNil: nil,
	}

	modelIn := userNew

	err := client.AllFieldTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, exists, err := client.AllFieldTypesRepo().Read(ctx, modelIn.ID())
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("model not found")
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)

	modelOut, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.Duration.Equal(time.Minute),
			where.AllFieldTypes.DurationPtr.GreaterThan(time.Minute),
			where.AllFieldTypes.DurationNil.Nil(true),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)
}

func TestUUID(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ptr := uuid.New()

	userNew := &model.AllFieldTypes{
		UUID:    uuid.New(),
		UUIDPtr: &ptr,
		UUIDNil: nil,
	}

	modelIn := userNew

	err := client.AllFieldTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, exists, err := client.AllFieldTypesRepo().Read(ctx, modelIn.ID())
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("model not found")
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)

	modelOut, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.UUID.Equal(modelIn.UUID),
			where.AllFieldTypes.UUIDPtr.Equal(*modelIn.UUIDPtr),
			where.AllFieldTypes.UUIDNil.Nil(true),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.ID{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)
}

func TestPassword(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	plainPassword := "test_password_123"

	// Step 1: Create a model with a known password (password is now in Login struct)
	modelIn := &model.AllFieldTypes{
		String: "password_test_user",
		Login: model.Login{
			Username: "testuser",
			Password: som.Password[som.Bcrypt](plainPassword),
		},
	}

	if string(modelIn.Login.Password) != plainPassword {
		t.Fatal("password should still be plaintext")
	}

	err := client.AllFieldTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	if string(modelIn.Login.Password) == plainPassword {
		t.Fatal("password should be hashed, not stored as plaintext")
	}

	// Step 2: Verify password was hashed (not equal to original)
	modelOut, exists, err := client.AllFieldTypesRepo().Read(ctx, modelIn.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("model not found")
	}

	if string(modelOut.Login.Password) == plainPassword {
		t.Fatal("password should be hashed, not stored as plaintext")
	}

	// Step 3: Verify password comparison works
	modelFound, err := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(modelIn.ID()),
			where.AllFieldTypes.Login().Password.Matches(plainPassword),
		).
		First(ctx)

	if err != nil {
		t.Fatalf("password comparison query failed: %v", err)
	}
	if modelFound == nil {
		t.Fatal("password comparison should have found the model")
	}

	// Step 4: Update OTHER field (not password)
	modelOut.Login.Username = "updated_user_name"

	err = client.AllFieldTypesRepo().Update(ctx, modelOut)
	if err != nil {
		t.Fatalf("failed to update model: %v", err)
	}

	// Step 5: Verify password comparison STILL works after update
	// This will FAIL if double-hashing occurs
	modelFoundAfterUpdate, err := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.ID.Equal(modelIn.ID()),
			where.AllFieldTypes.Login().Password.Matches(plainPassword),
		).
		First(ctx)

	if err != nil {
		t.Fatalf("password comparison after update failed: %v", err)
	}
	if modelFoundAfterUpdate == nil {
		t.Fatal("password comparison should still work after updating other fields - possible double-hashing issue")
	}

	assert.Equal(t, "updated_user_name", modelFoundAfterUpdate.Login.Username)
}

func FuzzWithDatabase(f *testing.F) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, f)
	defer cleanup()

	f.Add("Some User")

	f.Fuzz(func(t *testing.T, str string) {
		userIn := &model.AllFieldTypes{
			String: str,
		}

		err := client.AllFieldTypesRepo().Create(ctx, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == nil {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllFieldTypesRepo().Query().
			Filter(
				where.AllFieldTypes.ID.Equal(userIn.ID()),
			).
			First(ctx)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, userIn.String, userOut.String)
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

		userIn := &model.AllFieldTypes{
			String: "1",
		}

		err := client.AllFieldTypesRepo().CreateWithID(ctx, id, userIn)
		if err != nil {
			t.Fatal(err)
		}

		if userIn.ID() == nil {
			t.Fatal("user ID must not be empty after create call")
		}

		userOut, ok, err := client.AllFieldTypesRepo().Read(ctx, userIn.ID())

		if err != nil {
			t.Fatal(err)
		}

		if !ok {
			t.Fatalf("user with ID '%s' not found", userIn.ID())
		}

		assert.Equal(t, userIn.ID().String(), userOut.ID().String())
		assert.Equal(t, "1", userOut.String)

		userOut.String = "2"

		err = client.AllFieldTypesRepo().Update(ctx, userOut)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "2", userOut.String)

		err = client.AllFieldTypesRepo().Delete(ctx, userOut)
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
		userIn := &model.AllFieldTypes{
			String: "Some User",
		}

		err := client.AllFieldTypesRepo().Create(ctx, userIn)
		if err != nil {
			b.Fatal(err)
		}

		if userIn.ID() == nil {
			b.Fatal("user ID must not be empty after create call")
		}

		userOut, err := client.AllFieldTypesRepo().Query().
			Filter(
				where.AllFieldTypes.ID.Equal(userIn.ID()),
			).
			First(ctx)

		if err != nil {
			b.Fatal(err)
		}

		assert.Equal(b, userIn.String, userOut.String)
	}
}

func TestAsync(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{})
	if err != nil {
		t.Fatal(err)
	}

	resCh := client.AllFieldTypesRepo().Query().
		Filter().
		CountAsync(ctx)

	assert.NilError(t, <-resCh.Err())
	assert.Equal(t, 1, <-resCh.Val())

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{})
	if err != nil {
		t.Fatal(err)
	}

	resCh = client.AllFieldTypesRepo().Query().
		Filter().
		CountAsync(ctx)

	assert.NilError(t, <-resCh.Err())
	assert.Equal(t, 2, <-resCh.Val())
}

func TestRefresh(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	allFieldTypes := &model.AllFieldTypes{
		String: "some value",
	}

	err := client.AllFieldTypesRepo().Create(ctx, allFieldTypes)
	if err != nil {
		t.Fatal(err)
	}

	allFieldTypes.String = "some other value"

	err = client.AllFieldTypesRepo().Refresh(ctx, allFieldTypes)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "some value", allFieldTypes.String)
}
