package basic

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

var (
	uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	ulidRegex = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
)

func TestCustomID(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	t.Run("uuid", func(t *testing.T) {
		rec := model.SpecialTypes{Name: "uuid-test"}
		err := client.SpecialTypesRepo().Create(ctx, &rec)
		assert.NilError(t, err)
		assert.Assert(t, uuidRegex.MatchString(string(rec.ID())), "expected UUID format, got %q", rec.ID())
	})

	t.Run("rand", func(t *testing.T) {
		rec := model.SpecialRelation{Title: "rand-test"}
		err := client.SpecialRelationRepo().Create(ctx, &rec)
		assert.NilError(t, err)
		assert.Assert(t, rec.ID() != "", "expected non-empty ID")
	})

	t.Run("ulid", func(t *testing.T) {
		rec := model.AllTypes{FieldString: "ulid-test"}
		err := client.AllTypesRepo().Create(ctx, &rec)
		assert.NilError(t, err)
		assert.Assert(t, ulidRegex.MatchString(string(rec.ID())), "expected ULID format, got %q", rec.ID())
	})

	t.Run("create_with_id_uuid", func(t *testing.T) {
		knownUUID := "550e8400-e29b-41d4-a716-446655440000"
		rec := model.SpecialTypes{Name: "uuid-with-id"}
		err := client.SpecialTypesRepo().CreateWithID(ctx, knownUUID, &rec)
		assert.NilError(t, err)
		assert.Equal(t, string(rec.ID()), knownUUID)

		read, ok, err := client.SpecialTypesRepo().Read(ctx, knownUUID)
		assert.NilError(t, err)
		assert.Assert(t, ok, "expected record to exist")
		assert.Equal(t, string(read.ID()), knownUUID)
		assert.Equal(t, read.Name, "uuid-with-id")
	})

	t.Run("create_with_id_ulid", func(t *testing.T) {
		knownID := "my-custom-id"
		rec := model.AllTypes{FieldString: "ulid-with-id"}
		err := client.AllTypesRepo().CreateWithID(ctx, knownID, &rec)
		assert.NilError(t, err)
		assert.Equal(t, string(rec.ID()), knownID)

		read, ok, err := client.AllTypesRepo().Read(ctx, knownID)
		assert.NilError(t, err)
		assert.Assert(t, ok, "expected record to exist")
		assert.Equal(t, string(read.ID()), knownID)
		assert.Equal(t, read.FieldString, "ulid-with-id")
	})
}
