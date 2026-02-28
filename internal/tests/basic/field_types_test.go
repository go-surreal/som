package basic

import (
	"context"
	"math"
	"net/url"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	gofrsuuid "github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestNumerics(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "user"

	// MAX

	userMax := model.AllTypes{
		FieldString: str,
		FieldMonth:  time.January,
		FieldInt:    math.MaxInt,
		FieldInt8:   math.MaxInt8,
		FieldInt16:  math.MaxInt16,
		FieldInt32:  math.MaxInt32,
		FieldInt64:  math.MaxInt64,
		//Uint:    1, //math.MaxUint,
		FieldUint8:  math.MaxUint8,
		FieldUint16: math.MaxUint16,
		FieldUint32: math.MaxUint32,
		//Uint64:  1, //math.MaxUint64,
		//Uintptr: 1, //math.MaxUint64,
		FieldFloat32: math.MaxFloat32,
		FieldFloat64: math.MaxFloat64,
		FieldRune:    math.MaxInt32,
	}

	userIn := userMax

	err := client.AllTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.ID.Equal(string(userIn.ID())),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t,
		userMax, *userOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)

	// MIN

	userMin := model.AllTypes{
		FieldString: str,
		FieldMonth:  time.January,
		FieldInt:    math.MinInt,
		FieldInt8:   math.MinInt8,
		FieldInt16:  math.MinInt16,
		FieldInt32:  math.MinInt32,
		FieldInt64:  math.MinInt64,
		//Uint:    math.MaxUint,
		FieldUint8:  0,
		FieldUint16: 0,
		FieldUint32: 0,
		//Uint64:  0,
		//Uintptr: 0,
		FieldFloat32: -math.MaxFloat32,
		FieldFloat64: -math.MaxFloat64,
		FieldRune:    math.MinInt32,
	}

	userIn = userMin

	err = client.AllTypesRepo().Create(ctx, &userIn)
	if err != nil {
		t.Fatal(err)
	}

	userOut, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.ID.Equal(string(userIn.ID())),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t,
		userMin, *userOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)
}

func TestSlice(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// initial nil slice

	user := &model.AllTypes{FieldMonth: time.January}

	err := client.AllTypesRepo().Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.FieldNestedDataSlice == nil)

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice().IsEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.FieldNestedDataSlice == nil)

	// empty slice

	user.FieldNestedDataSlice = []model.NestedData{}

	assert.Check(t, user.FieldNestedDataSlice != nil)

	err = client.AllTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice().Empty(true),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.FieldNestedDataSlice != nil)

	// non-empty slice with actual data

	str1 := "hello"
	num1 := 42
	now1 := time.Now().Truncate(time.Microsecond).UTC()
	id1 := uuid.New()

	user.FieldNestedDataSlice = []model.NestedData{{
		StringPtr: &str1,
		IntPtr:    &num1,
		TimePtr:   &now1,
		UuidPtr:   &id1,
	}}

	err = client.AllTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity: %v", err)
	}

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice().NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.FieldNestedDataSlice) == 1)
	assert.Check(t, user.FieldNestedDataSlice[0].StringPtr != nil && *user.FieldNestedDataSlice[0].StringPtr == str1)
	assert.Check(t, user.FieldNestedDataSlice[0].IntPtr != nil && *user.FieldNestedDataSlice[0].IntPtr == num1)
	assert.Check(t, user.FieldNestedDataSlice[0].TimePtr != nil && user.FieldNestedDataSlice[0].TimePtr.Equal(now1))
	assert.Check(t, user.FieldNestedDataSlice[0].UuidPtr != nil && *user.FieldNestedDataSlice[0].UuidPtr == id1)

	// multiple elements

	str2 := "world"
	num2 := 99

	user.FieldNestedDataSlice = []model.NestedData{
		{StringPtr: &str1, IntPtr: &num1},
		{StringPtr: &str2, IntPtr: &num2},
	}

	err = client.AllTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity with multiple elements: %v", err)
	}

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice().NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.FieldNestedDataSlice) == 2)
	assert.Check(t, *user.FieldNestedDataSlice[0].StringPtr == str1)
	assert.Check(t, *user.FieldNestedDataSlice[0].IntPtr == num1)
	assert.Check(t, *user.FieldNestedDataSlice[1].StringPtr == str2)
	assert.Check(t, *user.FieldNestedDataSlice[1].IntPtr == num2)

	// sub-filter: single condition

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice(
				filter.NestedData.StringPtr.Equal(str1),
			).NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.FieldNestedDataSlice) == 2)

	// sub-filter: multiple conditions (AND)

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice(
				filter.NestedData.StringPtr.Equal(str1),
				filter.NestedData.IntPtr.Equal(&num1),
			).NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.FieldNestedDataSlice) == 2)

	// sub-filter: no match

	users, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice(
				filter.NestedData.StringPtr.Equal(str1),
				filter.NestedData.IntPtr.Equal(&num2),
			).NotEmpty(),
		).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(users) == 0)

	// sub-filter: without filters (same as calling without sub-filters)

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataSlice().NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.FieldNestedDataSlice) == 2)

	// test FieldNestedDataPtrSlice ([]*NestedData)

	user.FieldNestedDataPtrSlice = []*model.NestedData{
		{StringPtr: &str1, IntPtr: &num1},
		{StringPtr: &str2, IntPtr: &num2},
	}

	err = client.AllTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity with FieldNestedDataPtrSlice: %v", err)
	}

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataPtrSlice().NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(user.FieldNestedDataPtrSlice) == 2)
	assert.Check(t, user.FieldNestedDataPtrSlice[0] != nil && *user.FieldNestedDataPtrSlice[0].StringPtr == str1)
	assert.Check(t, user.FieldNestedDataPtrSlice[1] != nil && *user.FieldNestedDataPtrSlice[1].StringPtr == str2)

	// test FieldNestedDataPtrSlicePtr (*[]*NestedData)

	ptrSlice := []*model.NestedData{
		{StringPtr: &str1, IntPtr: &num1},
	}
	user.FieldNestedDataPtrSlicePtr = &ptrSlice

	err = client.AllTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity with FieldNestedDataPtrSlicePtr: %v", err)
	}

	user, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldNestedDataPtrSlicePtr().NotEmpty(),
		).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, user.FieldNestedDataPtrSlicePtr != nil)
	assert.Check(t, len(*user.FieldNestedDataPtrSlicePtr) == 1)
	assert.Check(t, (*user.FieldNestedDataPtrSlicePtr)[0] != nil && *(*user.FieldNestedDataPtrSlicePtr)[0].StringPtr == str1)

	// test refresh with struct slice data

	str3 := "refreshed"
	user.FieldNestedDataSlice = []model.NestedData{{StringPtr: &str3}}

	err = client.AllTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatalf("could not update entity for refresh test: %v", err)
	}

	// modify local data
	modified := "modified"
	user.FieldNestedDataSlice[0].StringPtr = &modified

	// refresh should restore to DB value
	err = client.AllTypesRepo().Refresh(ctx, user)
	if err != nil {
		t.Fatalf("could not refresh entity: %v", err)
	}

	assert.Check(t, len(user.FieldNestedDataSlice) == 1)
	assert.Check(t, *user.FieldNestedDataSlice[0].StringPtr == str3)
}

func TestSliceNilElements(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	s1 := "hello"
	s2 := "world"
	n1 := 10
	n2 := 20

	user := &model.AllTypes{
		FieldMonth:          time.January,
		FieldStringPtrSlice: []*string{&s1, nil, &s2},
		FieldIntPtrSlice:    []*int{&n1, nil, &n2},
	}

	err := client.AllTypesRepo().Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	readBack, _, err := client.AllTypesRepo().Read(ctx, string(user.ID()))
	if err != nil {
		t.Fatal(err)
	}

	// []*string with nil elements
	assert.Check(t, len(readBack.FieldStringPtrSlice) == 3)
	assert.Check(t, readBack.FieldStringPtrSlice[0] != nil && *readBack.FieldStringPtrSlice[0] == s1)
	assert.Check(t, readBack.FieldStringPtrSlice[1] == nil)
	assert.Check(t, readBack.FieldStringPtrSlice[2] != nil && *readBack.FieldStringPtrSlice[2] == s2)

	// []*int with nil elements
	assert.Check(t, len(readBack.FieldIntPtrSlice) == 3)
	assert.Check(t, readBack.FieldIntPtrSlice[0] != nil && *readBack.FieldIntPtrSlice[0] == n1)
	assert.Check(t, readBack.FieldIntPtrSlice[1] == nil)
	assert.Check(t, readBack.FieldIntPtrSlice[2] != nil && *readBack.FieldIntPtrSlice[2] == n2)

	// *[]*string with nil elements
	ptrSlice := []*string{&s1, nil, &s2}
	user2 := &model.AllTypes{
		FieldMonth:             time.January,
		FieldStringPtrSlicePtr: &ptrSlice,
	}

	err = client.AllTypesRepo().Create(ctx, user2)
	if err != nil {
		t.Fatal(err)
	}

	readBack2, _, err := client.AllTypesRepo().Read(ctx, string(user2.ID()))
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, readBack2.FieldStringPtrSlicePtr != nil)
	assert.Check(t, len(*readBack2.FieldStringPtrSlicePtr) == 3)
	assert.Check(t, (*readBack2.FieldStringPtrSlicePtr)[0] != nil && *(*readBack2.FieldStringPtrSlicePtr)[0] == s1)
	assert.Check(t, (*readBack2.FieldStringPtrSlicePtr)[1] == nil)
	assert.Check(t, (*readBack2.FieldStringPtrSlicePtr)[2] != nil && *(*readBack2.FieldStringPtrSlicePtr)[2] == s2)

	// []*NestedData with nil elements
	nested := model.NestedData{StringPtr: &s1}
	user3 := &model.AllTypes{
		FieldMonth:              time.January,
		FieldNestedDataPtrSlice: []*model.NestedData{&nested, nil},
	}

	err = client.AllTypesRepo().Create(ctx, user3)
	if err != nil {
		t.Fatal(err)
	}

	readBack3, _, err := client.AllTypesRepo().Read(ctx, string(user3.ID()))
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(readBack3.FieldNestedDataPtrSlice) == 2)
	assert.Check(t, readBack3.FieldNestedDataPtrSlice[0] != nil && *readBack3.FieldNestedDataPtrSlice[0].StringPtr == s1)
	assert.Check(t, readBack3.FieldNestedDataPtrSlice[1] == nil)

	// []*Role with nil elements
	r1 := model.RoleUser
	r2 := model.RoleAdmin
	user4 := &model.AllTypes{
		FieldMonth:        time.January,
		FieldEnumPtrSlice: []*model.Role{&r1, nil, &r2},
	}

	err = client.AllTypesRepo().Create(ctx, user4)
	if err != nil {
		t.Fatal(err)
	}

	readBack4, _, err := client.AllTypesRepo().Read(ctx, string(user4.ID()))
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(readBack4.FieldEnumPtrSlice) == 3)
	assert.Check(t, readBack4.FieldEnumPtrSlice[0] != nil && *readBack4.FieldEnumPtrSlice[0] == model.RoleUser)
	assert.Check(t, readBack4.FieldEnumPtrSlice[1] == nil)
	assert.Check(t, readBack4.FieldEnumPtrSlice[2] != nil && *readBack4.FieldEnumPtrSlice[2] == model.RoleAdmin)

	// []*SpecialTypes with nil elements
	node1 := &model.SpecialTypes{Name: "node1"}
	node2 := &model.SpecialTypes{Name: "node2"}

	err = client.SpecialTypesRepo().Create(ctx, node1)
	if err != nil {
		t.Fatal(err)
	}

	err = client.SpecialTypesRepo().Create(ctx, node2)
	if err != nil {
		t.Fatal(err)
	}

	user5 := &model.AllTypes{
		FieldMonth:        time.January,
		FieldNodePtrSlice: []*model.SpecialTypes{node1, nil, node2},
	}

	err = client.AllTypesRepo().Create(ctx, user5)
	if err != nil {
		t.Fatal(err)
	}

	readBack5, _, err := client.AllTypesRepo().Read(ctx, string(user5.ID()))
	if err != nil {
		t.Fatal(err)
	}

	// Nils are dropped for node pointer slices (by design).
	assert.Check(t, len(readBack5.FieldNodePtrSlice) == 2)
}

func TestTimestamps(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user := &model.AllTypes{FieldMonth: time.January}

	err := client.AllTypesRepo().Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, !user.CreatedAt().IsZero())
	assert.Check(t, !user.UpdatedAt().IsZero())
	assert.Check(t, time.Since(user.CreatedAt()) < 5*time.Second)
	assert.Check(t, time.Since(user.UpdatedAt()) < 5*time.Second)

	time.Sleep(time.Second)

	err = client.AllTypesRepo().Update(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, !user.CreatedAt().IsZero())
	assert.Check(t, !user.UpdatedAt().IsZero())
	assert.Check(t, time.Since(user.CreatedAt()) > time.Second)
	assert.Check(t, time.Since(user.UpdatedAt()) < 5*time.Second)
}

func TestURLTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	someURL, err := url.Parse("https://surrealdb.com")
	if err != nil {
		t.Fatal(err)
	}

	newModel := &model.AllTypes{
		FieldMonth:  time.January,
		FieldURLPtr: someURL,
		FieldURL:    *someURL,
	}

	err = client.AllTypesRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	readModel, exists, err := client.AllTypesRepo().Read(ctx, string(newModel.ID()))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, exists)

	assert.Equal(t, someURL.String(), readModel.FieldURLPtr.String())
	assert.Equal(t, someURL.String(), readModel.FieldURL.String())

	someURL, err = url.Parse("https://github.com/surrealdb/surrealdb")
	if err != nil {
		t.Fatal(err)
	}

	readModel.FieldURLPtr = someURL
	readModel.FieldURL = *someURL

	err = client.AllTypesRepo().Update(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, someURL.String(), readModel.FieldURLPtr.String())
	assert.Equal(t, someURL.String(), readModel.FieldURL.String())

	queryModel, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldURL.Equal(*someURL),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, someURL.String(), queryModel.FieldURLPtr.String())

	err = client.AllTypesRepo().Delete(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDuration(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ptr := time.Hour

	userNew := &model.AllTypes{
		FieldMonth:       time.January,
		FieldDuration:    time.Minute,
		FieldDurationPtr: &ptr,
		FieldDurationNil: nil,
	}

	modelIn := userNew

	err := client.AllTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, exists, err := client.AllTypesRepo().Read(ctx, string(modelIn.ID()))
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("model not found")
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)

	modelOut, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldDuration.Equal(time.Minute),
			filter.AllTypes.FieldDurationPtr.GreaterThan(time.Minute),
			filter.AllTypes.FieldDurationNil.Nil(true),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)
}

func TestUUID(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ptr := uuid.New()

	userNew := &model.AllTypes{
		FieldMonth:   time.January,
		FieldUUID:    uuid.New(),
		FieldUUIDPtr: &ptr,
		FieldUUIDNil: nil,
	}

	modelIn := userNew

	err := client.AllTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, exists, err := client.AllTypesRepo().Read(ctx, string(modelIn.ID()))
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("model not found")
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)

	modelOut, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldUUID.Equal(modelIn.FieldUUID),
			filter.AllTypes.FieldUUIDPtr.Equal(*modelIn.FieldUUIDPtr),
			filter.AllTypes.FieldUUIDNil.Nil(true),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)
}

func TestUUIDGofrs(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ptr := gofrsuuid.Must(gofrsuuid.NewV4())

	userNew := &model.AllTypes{
		FieldMonth:        time.January,
		FieldUUIDGofrs:    gofrsuuid.Must(gofrsuuid.NewV4()),
		FieldUUIDGofrsPtr: &ptr,
		FieldUUIDGofrsNil: nil,
	}

	modelIn := userNew

	err := client.AllTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, exists, err := client.AllTypesRepo().Read(ctx, string(modelIn.ID()))
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("model not found")
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)

	modelOut, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldUUIDGofrs.Equal(modelIn.FieldUUIDGofrs),
			filter.AllTypes.FieldUUIDGofrsPtr.Equal(*modelIn.FieldUUIDGofrsPtr),
			filter.AllTypes.FieldUUIDGofrsNil.Nil(true),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)
}

func TestPassword(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	plainPassword := "test_password_123"

	// Step 1: Create a model with a known password (password is now in Credentials struct)
	modelIn := &model.AllTypes{
		FieldString: "password_test_user",
		FieldMonth:  time.January,
		FieldCredentials: model.Credentials{
			Username: "testuser",
			Password: som.Password[som.Bcrypt](plainPassword),
		},
	}

	if string(modelIn.FieldCredentials.Password) != plainPassword {
		t.Fatal("password should still be plaintext")
	}

	err := client.AllTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	if string(modelIn.FieldCredentials.Password) == plainPassword {
		t.Fatal("password should be hashed, not stored as plaintext")
	}

	// Step 2: Verify password was hashed (not equal to original)
	modelOut, exists, err := client.AllTypesRepo().Read(ctx, string(modelIn.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("model not found")
	}

	if string(modelOut.FieldCredentials.Password) == plainPassword {
		t.Fatal("password should be hashed, not stored as plaintext")
	}

	// Step 3: Verify password comparison works
	modelFound, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.ID.Equal(string(modelIn.ID())),
			filter.AllTypes.FieldCredentials().Password.Verify(plainPassword),
		).
		First(ctx)

	if err != nil {
		t.Fatalf("password comparison query failed: %v", err)
	}
	if modelFound == nil {
		t.Fatal("password comparison should have found the model")
	}

	// Step 4: Update OTHER field (not password)
	modelOut.FieldCredentials.Username = "updated_user_name"

	err = client.AllTypesRepo().Update(ctx, modelOut)
	if err != nil {
		t.Fatalf("failed to update model: %v", err)
	}

	// Step 5: Verify password comparison STILL works after update
	// This will FAIL if double-hashing occurs
	modelFoundAfterUpdate, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.ID.Equal(string(modelIn.ID())),
			filter.AllTypes.FieldCredentials().Password.Verify(plainPassword),
		).
		First(ctx)

	if err != nil {
		t.Fatalf("password comparison after update failed: %v", err)
	}
	if modelFoundAfterUpdate == nil {
		t.Fatal("password comparison should still work after updating other fields - possible double-hashing issue")
	}

	assert.Equal(t, "updated_user_name", modelFoundAfterUpdate.FieldCredentials.Username)
}

func TestEmail(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	emailValue := som.Email("testuser@example.com")
	emailPtr := som.Email("admin@test.org")

	userNew := &model.AllTypes{
		FieldMonth:    time.January,
		FieldEmail:    emailValue,
		FieldEmailPtr: &emailPtr,
		FieldEmailNil: nil,
	}

	modelIn := userNew

	err := client.AllTypesRepo().Create(ctx, modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, exists, err := client.AllTypesRepo().Read(ctx, string(modelIn.ID()))
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("model not found")
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)

	modelOut, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldEmail.Equal(emailValue),
			filter.AllTypes.FieldEmailPtr.Equal(emailPtr),
			filter.AllTypes.FieldEmailNil.Nil(true),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, modelIn, modelOut,
		cmpopts.IgnoreUnexported(som.Node[som.ULID]{}, som.Node[som.UUID]{}, som.Timestamps{}, som.OptimisticLock{}, som.SoftDelete{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
		cmpopts.IgnoreFields(model.AllTypes{}, "FieldHookStatus", "FieldHookDetail"),
	)

	// Test email-specific filter methods
	modelOut, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldEmail.User().Equal("testuser"),
			filter.AllTypes.FieldEmail.Host().Equal("example.com"),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, emailValue, modelOut.FieldEmail)
}
