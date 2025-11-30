package basic

import (
	"context"
	"math"
	"net/url"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/by"
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
		"SELECT * FROM all_field_types WHERE (->group_member[WHERE (created_at < $A)]->group[WHERE (id = $B)] "+
			"AND duration::days(duration) < $C)",
		query.Describe(),
	)

	query = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StringPtr.Base64Decode().Base64Encode().
				Equal_(where.AllFieldTypes.String.Base64Decode().Base64Encode()),
		)

	assert.Equal(t,
		"SELECT * FROM all_field_types WHERE "+
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}),
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}),
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}),
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}, som.ID{}),
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}, som.ID{}),
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}, som.ID{}),
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}, som.ID{}),
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
			where.AllFieldTypes.Login().Password.Verify(plainPassword),
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
			where.AllFieldTypes.Login().Password.Verify(plainPassword),
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

func TestEmail(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	emailValue := som.Email("testuser@example.com")
	emailPtr := som.Email("admin@test.org")

	userNew := &model.AllFieldTypes{
		Email:    emailValue,
		EmailPtr: &emailPtr,
		EmailNil: nil,
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
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}, som.ID{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)

	modelOut, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.Email.Equal(emailValue),
			where.AllFieldTypes.EmailPtr.Equal(emailPtr),
			where.AllFieldTypes.EmailNil.Nil(true),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}, som.ID{}),
		cmpopts.IgnoreFields(model.Login{}, "Password", "PasswordPtr"),
	)

	// Test email-specific filter methods
	modelOut, err = client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.Email.User().Equal("testuser"),
			where.AllFieldTypes.Email.Host().Equal("example.com"),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, emailValue, modelOut.Email)
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

func TestSort(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	boolTrue := true
	boolFalse := false

	int8Val1 := int8(1)
	int8Val2 := int8(2)
	int8Val3 := int8(3)

	int16Val1 := int16(10)
	int16Val2 := int16(20)
	int16Val3 := int16(30)

	int32Val1 := int32(100)
	int32Val2 := int32(200)
	int32Val3 := int32(300)

	int64Val1 := int64(1000)
	int64Val2 := int64(2000)
	int64Val3 := int64(3000)

	uint8Val1 := uint8(1)
	uint8Val2 := uint8(2)
	uint8Val3 := uint8(3)

	uint16Val1 := uint16(10)
	uint16Val2 := uint16(20)
	uint16Val3 := uint16(30)

	uint32Val1 := uint32(100)
	uint32Val2 := uint32(200)
	uint32Val3 := uint32(300)

	strA := "A"
	strB := "B"
	strC := "C"

	now := time.Now().Truncate(time.Microsecond).UTC()
	time1 := now.Add(-2 * time.Hour)
	time2 := now.Add(-1 * time.Hour)
	time3 := now

	dur1 := time.Minute
	dur2 := time.Hour
	dur3 := 24 * time.Hour

	uuid1 := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	uuid2 := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	uuid3 := uuid.MustParse("00000000-0000-0000-0000-000000000003")

	url1, _ := url.Parse("https://a.com")
	url2, _ := url.Parse("https://b.com")
	url3, _ := url.Parse("https://c.com")

	record1 := &model.AllFieldTypes{
		String:    "charlie",
		StringPtr: &strC,
		Int:       30,
		IntPtr:    func() *int { v := 300; return &v }(),
		Int8:      3,
		Int8Ptr:   &int8Val3,
		Int16:     30,
		Int16Ptr:  &int16Val3,
		Int32:     300,
		Int32Ptr:  &int32Val3,
		Int64:     3000,
		Int64Ptr:  &int64Val3,
		Uint8:     3,
		Uint8Ptr:  &uint8Val3,
		Uint16:    30,
		Uint16Ptr: &uint16Val3,
		Uint32:    300,
		Uint32Ptr: &uint32Val3,
		Float32:   3.0,
		Float64:   30.0,
		Rune:      'c',
		Bool:      true,
		BoolPtr:   &boolTrue,
		Byte:      3,
		BytePtr:   &uint8Val3,
		Time:      time3,
		TimePtr:   &time3,
		Duration:  dur3,
		UUID:      uuid3,
		UUIDPtr:   &uuid3,
		URL:       *url3,
		URLPtr:    url3,
		Role:      model.RoleUser,
		Login:     model.Login{Username: "charlie", Password: "pass3"},
	}

	record2 := &model.AllFieldTypes{
		String:    "alpha",
		StringPtr: &strA,
		Int:       10,
		IntPtr:    func() *int { v := 100; return &v }(),
		Int8:      1,
		Int8Ptr:   &int8Val1,
		Int16:     10,
		Int16Ptr:  &int16Val1,
		Int32:     100,
		Int32Ptr:  &int32Val1,
		Int64:     1000,
		Int64Ptr:  &int64Val1,
		Uint8:     1,
		Uint8Ptr:  &uint8Val1,
		Uint16:    10,
		Uint16Ptr: &uint16Val1,
		Uint32:    100,
		Uint32Ptr: &uint32Val1,
		Float32:   1.0,
		Float64:   10.0,
		Rune:      'a',
		Bool:      false,
		BoolPtr:   &boolFalse,
		Byte:      1,
		BytePtr:   &uint8Val1,
		Time:      time1,
		TimePtr:   &time1,
		Duration:  dur1,
		UUID:      uuid1,
		UUIDPtr:   &uuid1,
		URL:       *url1,
		URLPtr:    url1,
		Role:      model.RoleAdmin,
		Login:     model.Login{Username: "alpha", Password: "pass1"},
	}

	record3 := &model.AllFieldTypes{
		String:    "bravo",
		StringPtr: &strB,
		Int:       20,
		IntPtr:    func() *int { v := 200; return &v }(),
		Int8:      2,
		Int8Ptr:   &int8Val2,
		Int16:     20,
		Int16Ptr:  &int16Val2,
		Int32:     200,
		Int32Ptr:  &int32Val2,
		Int64:     2000,
		Int64Ptr:  &int64Val2,
		Uint8:     2,
		Uint8Ptr:  &uint8Val2,
		Uint16:    20,
		Uint16Ptr: &uint16Val2,
		Uint32:    200,
		Uint32Ptr: &uint32Val2,
		Float32:   2.0,
		Float64:   20.0,
		Rune:      'b',
		Bool:      true,
		BoolPtr:   &boolTrue,
		Byte:      2,
		BytePtr:   &uint8Val2,
		Time:      time2,
		TimePtr:   &time2,
		Duration:  dur2,
		UUID:      uuid2,
		UUIDPtr:   &uuid2,
		URL:       *url2,
		URLPtr:    url2,
		Role:      model.RoleUser,
		Login:     model.Login{Username: "bravo", Password: "pass2"},
	}

	for _, r := range []*model.AllFieldTypes{record1, record2, record3} {
		if err := client.AllFieldTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("String", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.String.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, "alpha", results[0].String)
		assert.Equal(t, "bravo", results[1].String)
		assert.Equal(t, "charlie", results[2].String)

		results, err = client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.String.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "charlie", results[0].String)
		assert.Equal(t, "bravo", results[1].String)
		assert.Equal(t, "alpha", results[2].String)
	})

	t.Run("StringCollate", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.String.Collate().Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, "alpha", results[0].String)
	})

	t.Run("StringNumeric", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.String.Numeric().Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
	})

	t.Run("Int", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Int.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 10, results[0].Int)
		assert.Equal(t, 20, results[1].Int)
		assert.Equal(t, 30, results[2].Int)

		results, err = client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Int.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 30, results[0].Int)
		assert.Equal(t, 20, results[1].Int)
		assert.Equal(t, 10, results[2].Int)
	})

	t.Run("Int8", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Int8.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int8(1), results[0].Int8)
		assert.Equal(t, int8(2), results[1].Int8)
		assert.Equal(t, int8(3), results[2].Int8)
	})

	t.Run("Int16", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Int16.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int16(10), results[0].Int16)
		assert.Equal(t, int16(20), results[1].Int16)
		assert.Equal(t, int16(30), results[2].Int16)
	})

	t.Run("Int32", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Int32.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int32(100), results[0].Int32)
		assert.Equal(t, int32(200), results[1].Int32)
		assert.Equal(t, int32(300), results[2].Int32)
	})

	t.Run("Int64", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Int64.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int64(1000), results[0].Int64)
		assert.Equal(t, int64(2000), results[1].Int64)
		assert.Equal(t, int64(3000), results[2].Int64)
	})

	t.Run("Uint8", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Uint8.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint8(1), results[0].Uint8)
		assert.Equal(t, uint8(2), results[1].Uint8)
		assert.Equal(t, uint8(3), results[2].Uint8)
	})

	t.Run("Uint16", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Uint16.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint16(10), results[0].Uint16)
		assert.Equal(t, uint16(20), results[1].Uint16)
		assert.Equal(t, uint16(30), results[2].Uint16)
	})

	t.Run("Uint32", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Uint32.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint32(100), results[0].Uint32)
		assert.Equal(t, uint32(200), results[1].Uint32)
		assert.Equal(t, uint32(300), results[2].Uint32)
	})

	t.Run("Float32", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Float32.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, float32(1.0), results[0].Float32)
		assert.Equal(t, float32(2.0), results[1].Float32)
		assert.Equal(t, float32(3.0), results[2].Float32)
	})

	t.Run("Float64", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Float64.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 10.0, results[0].Float64)
		assert.Equal(t, 20.0, results[1].Float64)
		assert.Equal(t, 30.0, results[2].Float64)
	})

	t.Run("Rune", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Rune.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 'a', results[0].Rune)
		assert.Equal(t, 'b', results[1].Rune)
		assert.Equal(t, 'c', results[2].Rune)
	})

	t.Run("Bool", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Bool.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, false, results[0].Bool)

		results, err = client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Bool.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, true, results[0].Bool)
	})

	t.Run("Byte", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Byte.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, byte(1), results[0].Byte)
		assert.Equal(t, byte(2), results[1].Byte)
		assert.Equal(t, byte(3), results[2].Byte)
	})

	t.Run("Time", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Time.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Check(t, results[0].Time.Before(results[1].Time))
		assert.Check(t, results[1].Time.Before(results[2].Time))

		results, err = client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Time.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Check(t, results[0].Time.After(results[1].Time))
		assert.Check(t, results[1].Time.After(results[2].Time))
	})

	t.Run("Duration", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Duration.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, time.Minute, results[0].Duration)
		assert.Equal(t, time.Hour, results[1].Duration)
		assert.Equal(t, 24*time.Hour, results[2].Duration)
	})

	t.Run("UUID", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.UUID.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uuid1, results[0].UUID)
		assert.Equal(t, uuid2, results[1].UUID)
		assert.Equal(t, uuid3, results[2].UUID)
	})

	t.Run("URL", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.URL.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "https://a.com", results[0].URL.String())
		assert.Equal(t, "https://b.com", results[1].URL.String())
		assert.Equal(t, "https://c.com", results[2].URL.String())
	})

	t.Run("Enum", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Role.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, model.RoleAdmin, results[0].Role)
	})

	t.Run("NestedStruct", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Login().Username.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "alpha", results[0].Login.Username)
		assert.Equal(t, "bravo", results[1].Login.Username)
		assert.Equal(t, "charlie", results[2].Login.Username)

		results, err = client.AllFieldTypesRepo().Query().
			Order(by.AllFieldTypes.Login().Username.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "charlie", results[0].Login.Username)
		assert.Equal(t, "bravo", results[1].Login.Username)
		assert.Equal(t, "alpha", results[2].Login.Username)
	})

	t.Run("MultipleFields", func(t *testing.T) {
		results, err := client.AllFieldTypesRepo().Query().
			Order(
				by.AllFieldTypes.Bool.Asc(),
				by.AllFieldTypes.String.Asc(),
			).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, false, results[0].Bool)
		assert.Equal(t, "alpha", results[0].String)
	})
}
