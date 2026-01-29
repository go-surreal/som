package basic

import (
	"context"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/field"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestDistinct(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("EmptyTable", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldString)
		assert.NilError(t, err)
		assert.Check(t, vals == nil)
	})

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

	url1, _ := url.Parse("https://example.com/alpha")
	url2, _ := url.Parse("https://example.com/bravo")
	url3, _ := url.Parse("https://example.com/charlie")

	record1 := &model.AllTypes{
		FieldString:      "alpha",
		FieldInt:         10,
		FieldBool:        false,
		FieldTime:        time1,
		FieldDuration:    dur1,
		FieldUUID:        uuid1,
		FieldURL:         *url1,
		FieldEnum:        model.RoleAdmin,
		FieldCredentials: model.Credentials{Username: "alice", Password: "pass1"},
	}

	record2 := &model.AllTypes{
		FieldString:      "bravo",
		FieldInt:         20,
		FieldBool:        true,
		FieldTime:        time2,
		FieldDuration:    dur2,
		FieldUUID:        uuid2,
		FieldURL:         *url2,
		FieldEnum:        model.RoleUser,
		FieldCredentials: model.Credentials{Username: "bob", Password: "pass2"},
	}

	record3 := &model.AllTypes{
		FieldString:      "charlie",
		FieldInt:         30,
		FieldBool:        true,
		FieldTime:        time3,
		FieldDuration:    dur3,
		FieldUUID:        uuid3,
		FieldURL:         *url3,
		FieldEnum:        model.RoleUser,
		FieldCredentials: model.Credentials{Username: "charlie", Password: "pass3"},
	}

	for _, r := range []*model.AllTypes{record1, record2, record3} {
		if err := client.AllTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("String", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldString)
		assert.NilError(t, err)
		sort.Strings(vals)
		assert.DeepEqual(t, []string{"alpha", "bravo", "charlie"}, vals)
	})

	t.Run("Int", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldInt)
		assert.NilError(t, err)
		sort.Ints(vals)
		assert.DeepEqual(t, []int{10, 20, 30}, vals)
	})

	t.Run("Bool", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldBool)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return !vals[i] && vals[j] })
		assert.DeepEqual(t, []bool{false, true}, vals)
	})

	t.Run("Enum", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldEnum)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		assert.DeepEqual(t, []model.Role{model.RoleAdmin, model.RoleUser}, vals)
	})

	t.Run("Time", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldTime)
		assert.NilError(t, err)
		assert.Equal(t, 3, len(vals))
	})

	t.Run("Duration", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldDuration)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		assert.DeepEqual(t, []time.Duration{dur1, dur2, dur3}, vals)
	})

	t.Run("UUID", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldUUID)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i].String() < vals[j].String() })
		assert.DeepEqual(t, []uuid.UUID{uuid1, uuid2, uuid3}, vals)
	})

	t.Run("URL", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldURL)
		assert.NilError(t, err)
		sort.Slice(vals, func(i, j int) bool { return vals[i].String() < vals[j].String() })
		assert.DeepEqual(t, []url.URL{*url1, *url2, *url3}, vals)
	})

	t.Run("NestedField", func(t *testing.T) {
		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldCredentials().Username)
		assert.NilError(t, err)
		sort.Strings(vals)
		assert.DeepEqual(t, []string{"alice", "bob", "charlie"}, vals)
	})

	t.Run("WithFilter", func(t *testing.T) {
		vals, err := query.Distinct(ctx,
			client.AllTypesRepo().Query().Where(
				filter.AllTypes.FieldEnum.Equal(model.RoleAdmin),
			),
			field.AllTypes.FieldString,
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, []string{"alpha"}, vals)
	})

	t.Run("StringWithDuplicates", func(t *testing.T) {
		dup1 := &model.AllTypes{FieldString: "alpha", FieldCredentials: model.Credentials{Username: "dup1", Password: "pass"}}
		dup2 := &model.AllTypes{FieldString: "alpha", FieldCredentials: model.Credentials{Username: "dup2", Password: "pass"}}
		for _, r := range []*model.AllTypes{dup1, dup2} {
			if err := client.AllTypesRepo().Create(ctx, r); err != nil {
				t.Fatal(err)
			}
		}

		vals, err := query.Distinct(ctx, client.AllTypesRepo().Query(), field.AllTypes.FieldString)
		assert.NilError(t, err)
		sort.Strings(vals)
		assert.DeepEqual(t, []string{"alpha", "bravo", "charlie"}, vals)
	})
}
