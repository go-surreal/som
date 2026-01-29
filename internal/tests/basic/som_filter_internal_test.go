package basic

import (
	"context"
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestFilterCompareFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	str := "Some Value"
	date := time.Now()

	modelNew := model.AllTypes{
		FieldString:    str,
		FieldStringPtr: &str,

		FieldTime:    date.Add(-time.Hour),
		FieldTimePtr: &date,

		FieldDuration: time.Hour,
	}

	modelIn := modelNew

	err := client.AllTypesRepo().Create(ctx, &modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldStringPtr.Equal_(filter.AllTypes.FieldString),
			filter.AllTypes.FieldTimePtr.After_(filter.AllTypes.FieldTime),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, str, modelOut.FieldString)
	assert.Equal(t, str, *modelOut.FieldStringPtr)

	assert.DeepEqual(t,
		modelNew, *modelOut,
		cmpopts.IgnoreUnexported(som.Node{}, som.Timestamps{}, som.OptimisticLock{}),
		cmpopts.IgnoreFields(model.Credentials{}, "Password", "PasswordPtr"),
	)
}
