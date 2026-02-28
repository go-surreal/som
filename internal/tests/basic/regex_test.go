package basic

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestRegex(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	pattern := regexp.MustCompile(`^[a-z]+$`)

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldMonth: time.January,
		FieldRegex: *pattern,
	})
	if err != nil {
		t.Fatalf("create with regex failed: %v", err)
	}

	res, err := client.AllTypesRepo().Query().
		Where().
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(res))
	t.Logf("FieldRegex pattern: %q", res[0].FieldRegex.String())
	assert.Equal(t, `^[a-z]+$`, res[0].FieldRegex.String())
}
