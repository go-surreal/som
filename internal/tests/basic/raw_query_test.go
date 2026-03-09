package basic

import (
	"context"
	"errors"
	"testing"
	"time"

	"som.test/gen/som"
	"som.test/model"
	"gotest.tools/v3/assert"
)

func TestRawQuery_Scan(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 3; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: "raw_test",
			FieldInt:    i,
			FieldMonth:  time.January,
		})
		assert.NilError(t, err)
	}

	result, err := client.Raw(ctx, "SELECT * FROM all_types WHERE field_string = $name", som.Params{"name": "raw_test"})
	assert.NilError(t, err)

	var rows []map[string]any
	err = result.Scan(&rows)
	assert.NilError(t, err)
	assert.Equal(t, 3, len(rows))
}

func TestRawQuery_ScanOne(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "scan_one_test",
		FieldInt:    42,
		FieldMonth:  time.January,
	})
	assert.NilError(t, err)

	result, err := client.Raw(ctx, "SELECT * FROM all_types WHERE field_string = $name", som.Params{"name": "scan_one_test"})
	assert.NilError(t, err)

	var row map[string]any
	err = result.ScanOne(&row)
	assert.NilError(t, err)
	assert.Equal(t, "scan_one_test", row["field_string"])
}

func TestRawQuery_ScanOne_NotFound(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	result, err := client.Raw(ctx, "SELECT * FROM all_types WHERE field_string = 'does_not_exist'", nil)
	assert.NilError(t, err)

	var row map[string]any
	err = result.ScanOne(&row)
	assert.Assert(t, errors.Is(err, som.ErrNotFound))
}

func TestRawQuery_NoParams(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	result, err := client.Raw(ctx, "SELECT * FROM all_types LIMIT 0", nil)
	assert.NilError(t, err)

	var rows []map[string]any
	err = result.Scan(&rows)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(rows))
}

func TestRawQuery_DDL(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	_, err := client.Raw(ctx, "DEFINE TABLE IF NOT EXISTS raw_test_temp", nil)
	assert.NilError(t, err)

	_, err = client.Raw(ctx, "REMOVE TABLE IF EXISTS raw_test_temp", nil)
	assert.NilError(t, err)
}

func TestRawQuery_ParamBinding(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "param_test",
		FieldInt:    99,
		FieldMonth:  time.January,
	})
	assert.NilError(t, err)

	result, err := client.Raw(ctx,
		"SELECT * FROM all_types WHERE field_string = $name AND field_int = $val",
		som.Params{"name": "param_test", "val": 99},
	)
	assert.NilError(t, err)

	var rows []map[string]any
	err = result.Scan(&rows)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(rows))
}

func TestRawQuery_InvalidQuery(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	_, err := client.Raw(ctx, "THIS IS NOT VALID SURREALQL", nil)
	assert.Assert(t, err != nil)
}
