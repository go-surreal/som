package basic

import (
	"context"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/gen/som/repo"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestQuery(t *testing.T) {
	client := &repo.ClientImpl{}

	query := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.
				FieldEdgeRelations(
					filter.EdgeRelation.CreatedAt.Before(time.Now()),
				).
				SpecialTypes(
					filter.SpecialTypes.ID.Equal("some_id"),
				),

			filter.AllTypes.FieldDuration.Days().LessThan(4),

			//filter.AllTypes.Float64.Equal_(constant.E[model.AllTypes]()),
			//
			//constant.String[model.AllTypes]("A").Equal_(constant.String[model.AllTypes]("A")),
		)

	assert.Equal(t,
		"SELECT * FROM all_types WHERE (->edge_relation[WHERE (created_at < $A)]->special_types[WHERE (id = $B)] "+
			"AND duration::days(field_duration) < $C)",
		query.Describe(),
	)

	query = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldStringPtr.Base64Decode().Base64Encode().
				Equal_(filter.AllTypes.FieldString.Base64Decode().Base64Encode()),
		)

	assert.Equal(t,
		"SELECT * FROM all_types WHERE "+
			"(encoding::base64::encode(encoding::base64::decode(field_string_ptr)) "+
			"= encoding::base64::encode(encoding::base64::decode(field_string)))",
		query.Describe(),
	)
}

func TestQueryLimitOffset(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldInt: i,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	results, err := client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldInt.Asc()).
		Limit(2).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 0, results[0].FieldInt)
	assert.Equal(t, 1, results[1].FieldInt)

	results, err = client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldInt.Asc()).
		Start(2).
		Limit(2).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(results))
	assert.Equal(t, 2, results[0].FieldInt)
	assert.Equal(t, 3, results[1].FieldInt)

	results, err = client.AllTypesRepo().Query().
		Order(by.AllTypes.FieldInt.Asc()).
		Start(4).
		Limit(10).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 4, results[0].FieldInt)
}

func TestQueryIDs(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 3; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldInt: i,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	ids, err := client.AllTypesRepo().Query().AllIDs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(ids))
	for _, id := range ids {
		assert.Check(t, id != "")
	}

	id, err := client.AllTypesRepo().Query().FirstID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Check(t, id != "")
}

func TestOrderRandom(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 3; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldInt: i,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	results, err := client.AllTypesRepo().Query().
		OrderRandom().
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(results))
}
