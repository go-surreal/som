package basic

import (
	"context"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/constant"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestConstants(t *testing.T) {
	ctx := context.Background()
	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{FieldFloat64: 2.0})
	if err != nil {
		t.Fatal(err)
	}

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{FieldFloat64: 4.0})
	if err != nil {
		t.Fatal(err)
	}

	// Filter where FieldFloat64 < PI (~3.14)
	results, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldFloat64.LessThan_(constant.PI[model.AllTypes]())).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 2.0, results[0].FieldFloat64)

	// Filter where FieldFloat64 > E (~2.718)
	results, err = client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldFloat64.GreaterThan_(constant.E[model.AllTypes]())).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 4.0, results[0].FieldFloat64)

	// String constant: filter where FieldString equals a constant string
	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{FieldString: "hello"})
	if err != nil {
		t.Fatal(err)
	}

	results, err = client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldString.Equal_(constant.String[model.AllTypes]("hello"))).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "hello", results[0].FieldString)
}
