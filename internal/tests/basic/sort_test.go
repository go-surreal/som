package basic

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

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

	record1 := &model.AllTypes{
		FieldString:    "charlie",
		FieldStringPtr: &strC,
		FieldInt:       30,
		FieldIntPtr:    func() *int { v := 300; return &v }(),
		FieldInt8:      3,
		FieldInt8Ptr:   &int8Val3,
		FieldInt16:     30,
		FieldInt16Ptr:  &int16Val3,
		FieldInt32:     300,
		FieldInt32Ptr:  &int32Val3,
		FieldInt64:     3000,
		FieldInt64Ptr:  &int64Val3,
		FieldUint8:     3,
		FieldUint8Ptr:  &uint8Val3,
		FieldUint16:    30,
		FieldUint16Ptr: &uint16Val3,
		FieldUint32:    300,
		FieldUint32Ptr: &uint32Val3,
		FieldFloat32:   3.0,
		FieldFloat64:   30.0,
		FieldRune:      'c',
		FieldBool:      true,
		FieldBoolPtr:   &boolTrue,
		FieldByte:      3,
		FieldBytePtr:   &uint8Val3,
		FieldTime:      time3,
		FieldTimePtr:   &time3,
		FieldDuration:  dur3,
		FieldUUID:      uuid3,
		FieldUUIDPtr:   &uuid3,
		FieldURL:       *url3,
		FieldURLPtr:    url3,
		FieldEnum:      model.RoleUser,
		FieldCredentials: model.Credentials{Username: "charlie", Password: "pass3"},
		FieldMonth:       time.January,
	}

	record2 := &model.AllTypes{
		FieldString:    "alpha",
		FieldStringPtr: &strA,
		FieldInt:       10,
		FieldIntPtr:    func() *int { v := 100; return &v }(),
		FieldInt8:      1,
		FieldInt8Ptr:   &int8Val1,
		FieldInt16:     10,
		FieldInt16Ptr:  &int16Val1,
		FieldInt32:     100,
		FieldInt32Ptr:  &int32Val1,
		FieldInt64:     1000,
		FieldInt64Ptr:  &int64Val1,
		FieldUint8:     1,
		FieldUint8Ptr:  &uint8Val1,
		FieldUint16:    10,
		FieldUint16Ptr: &uint16Val1,
		FieldUint32:    100,
		FieldUint32Ptr: &uint32Val1,
		FieldFloat32:   1.0,
		FieldFloat64:   10.0,
		FieldRune:      'a',
		FieldBool:      false,
		FieldBoolPtr:   &boolFalse,
		FieldByte:      1,
		FieldBytePtr:   &uint8Val1,
		FieldTime:      time1,
		FieldTimePtr:   &time1,
		FieldDuration:  dur1,
		FieldUUID:      uuid1,
		FieldUUIDPtr:   &uuid1,
		FieldURL:       *url1,
		FieldURLPtr:    url1,
		FieldEnum:      model.RoleAdmin,
		FieldCredentials: model.Credentials{Username: "alpha", Password: "pass1"},
		FieldMonth:       time.January,
	}

	record3 := &model.AllTypes{
		FieldString:    "bravo",
		FieldStringPtr: &strB,
		FieldInt:       20,
		FieldIntPtr:    func() *int { v := 200; return &v }(),
		FieldInt8:      2,
		FieldInt8Ptr:   &int8Val2,
		FieldInt16:     20,
		FieldInt16Ptr:  &int16Val2,
		FieldInt32:     200,
		FieldInt32Ptr:  &int32Val2,
		FieldInt64:     2000,
		FieldInt64Ptr:  &int64Val2,
		FieldUint8:     2,
		FieldUint8Ptr:  &uint8Val2,
		FieldUint16:    20,
		FieldUint16Ptr: &uint16Val2,
		FieldUint32:    200,
		FieldUint32Ptr: &uint32Val2,
		FieldFloat32:   2.0,
		FieldFloat64:   20.0,
		FieldRune:      'b',
		FieldBool:      true,
		FieldBoolPtr:   &boolTrue,
		FieldByte:      2,
		FieldBytePtr:   &uint8Val2,
		FieldTime:      time2,
		FieldTimePtr:   &time2,
		FieldDuration:  dur2,
		FieldUUID:      uuid2,
		FieldUUIDPtr:   &uuid2,
		FieldURL:       *url2,
		FieldURLPtr:    url2,
		FieldEnum:      model.RoleUser,
		FieldCredentials: model.Credentials{Username: "bravo", Password: "pass2"},
		FieldMonth:       time.January,
	}

	for _, r := range []*model.AllTypes{record1, record2, record3} {
		if err := client.AllTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("String", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldString.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, "alpha", results[0].FieldString)
		assert.Equal(t, "bravo", results[1].FieldString)
		assert.Equal(t, "charlie", results[2].FieldString)

		results, err = client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldString.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "charlie", results[0].FieldString)
		assert.Equal(t, "bravo", results[1].FieldString)
		assert.Equal(t, "alpha", results[2].FieldString)
	})

	t.Run("StringCollate", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldString.Collate().Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, "alpha", results[0].FieldString)
	})

	t.Run("StringNumeric", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldString.Numeric().Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
	})

	t.Run("Int", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldInt.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 10, results[0].FieldInt)
		assert.Equal(t, 20, results[1].FieldInt)
		assert.Equal(t, 30, results[2].FieldInt)

		results, err = client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldInt.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 30, results[0].FieldInt)
		assert.Equal(t, 20, results[1].FieldInt)
		assert.Equal(t, 10, results[2].FieldInt)
	})

	t.Run("Int8", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldInt8.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int8(1), results[0].FieldInt8)
		assert.Equal(t, int8(2), results[1].FieldInt8)
		assert.Equal(t, int8(3), results[2].FieldInt8)
	})

	t.Run("Int16", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldInt16.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int16(10), results[0].FieldInt16)
		assert.Equal(t, int16(20), results[1].FieldInt16)
		assert.Equal(t, int16(30), results[2].FieldInt16)
	})

	t.Run("Int32", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldInt32.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int32(100), results[0].FieldInt32)
		assert.Equal(t, int32(200), results[1].FieldInt32)
		assert.Equal(t, int32(300), results[2].FieldInt32)
	})

	t.Run("Int64", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldInt64.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, int64(1000), results[0].FieldInt64)
		assert.Equal(t, int64(2000), results[1].FieldInt64)
		assert.Equal(t, int64(3000), results[2].FieldInt64)
	})

	t.Run("Uint8", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldUint8.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint8(1), results[0].FieldUint8)
		assert.Equal(t, uint8(2), results[1].FieldUint8)
		assert.Equal(t, uint8(3), results[2].FieldUint8)
	})

	t.Run("Uint16", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldUint16.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint16(10), results[0].FieldUint16)
		assert.Equal(t, uint16(20), results[1].FieldUint16)
		assert.Equal(t, uint16(30), results[2].FieldUint16)
	})

	t.Run("Uint32", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldUint32.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint32(100), results[0].FieldUint32)
		assert.Equal(t, uint32(200), results[1].FieldUint32)
		assert.Equal(t, uint32(300), results[2].FieldUint32)
	})

	t.Run("Float32", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldFloat32.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, float32(1.0), results[0].FieldFloat32)
		assert.Equal(t, float32(2.0), results[1].FieldFloat32)
		assert.Equal(t, float32(3.0), results[2].FieldFloat32)
	})

	t.Run("Float64", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldFloat64.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 10.0, results[0].FieldFloat64)
		assert.Equal(t, 20.0, results[1].FieldFloat64)
		assert.Equal(t, 30.0, results[2].FieldFloat64)
	})

	t.Run("Rune", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldRune.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 'a', results[0].FieldRune)
		assert.Equal(t, 'b', results[1].FieldRune)
		assert.Equal(t, 'c', results[2].FieldRune)
	})

	t.Run("Bool", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldBool.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, false, results[0].FieldBool)

		results, err = client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldBool.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, true, results[0].FieldBool)
	})

	t.Run("Byte", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldByte.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, byte(1), results[0].FieldByte)
		assert.Equal(t, byte(2), results[1].FieldByte)
		assert.Equal(t, byte(3), results[2].FieldByte)
	})

	t.Run("Time", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldTime.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Check(t, results[0].FieldTime.Before(results[1].FieldTime))
		assert.Check(t, results[1].FieldTime.Before(results[2].FieldTime))

		results, err = client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldTime.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Check(t, results[0].FieldTime.After(results[1].FieldTime))
		assert.Check(t, results[1].FieldTime.After(results[2].FieldTime))
	})

	t.Run("Duration", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldDuration.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, time.Minute, results[0].FieldDuration)
		assert.Equal(t, time.Hour, results[1].FieldDuration)
		assert.Equal(t, 24*time.Hour, results[2].FieldDuration)
	})

	t.Run("UUID", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldUUID.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uuid1, results[0].FieldUUID)
		assert.Equal(t, uuid2, results[1].FieldUUID)
		assert.Equal(t, uuid3, results[2].FieldUUID)
	})

	t.Run("URL", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldURL.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "https://a.com", results[0].FieldURL.String())
		assert.Equal(t, "https://b.com", results[1].FieldURL.String())
		assert.Equal(t, "https://c.com", results[2].FieldURL.String())
	})

	t.Run("Enum", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldEnum.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, model.RoleAdmin, results[0].FieldEnum)
	})

	t.Run("NestedStruct", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldCredentials().Username.Asc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "alpha", results[0].FieldCredentials.Username)
		assert.Equal(t, "bravo", results[1].FieldCredentials.Username)
		assert.Equal(t, "charlie", results[2].FieldCredentials.Username)

		results, err = client.AllTypesRepo().Query().
			Order(by.AllTypes.FieldCredentials().Username.Desc()).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "charlie", results[0].FieldCredentials.Username)
		assert.Equal(t, "bravo", results[1].FieldCredentials.Username)
		assert.Equal(t, "alpha", results[2].FieldCredentials.Username)
	})

	t.Run("MultipleFields", func(t *testing.T) {
		results, err := client.AllTypesRepo().Query().
			Order(
				by.AllTypes.FieldBool.Asc(),
				by.AllTypes.FieldString.Asc(),
			).All(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 3, len(results))
		assert.Equal(t, false, results[0].FieldBool)
		assert.Equal(t, "alpha", results[0].FieldString)
	})
}
