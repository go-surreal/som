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
		assert.Assert(t, uuidRegex.MatchString(rec.ID()), "expected UUID format, got %q", rec.ID())
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
		assert.Assert(t, ulidRegex.MatchString(rec.ID()), "expected ULID format, got %q", rec.ID())
	})
}
