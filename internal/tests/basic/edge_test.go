package basic

import (
	"context"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestEdgeRelation(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	allTypesNode := &model.AllTypes{
		FieldString: "edge_source",
	}
	err := client.AllTypesRepo().Create(ctx, allTypesNode)
	if err != nil {
		t.Fatal(err)
	}

	specialTypesNode := &model.SpecialTypes{
		Name: "edge_target",
	}
	err = client.SpecialTypesRepo().Create(ctx, specialTypesNode)
	if err != nil {
		t.Fatal(err)
	}

	edge := &model.EdgeRelation{
		AllTypes:     *allTypesNode,
		SpecialTypes: *specialTypesNode,
		Meta: model.EdgeMeta{
			IsAdmin:  true,
			IsActive: true,
		},
	}

	err = client.AllTypesRepo().Relate().FieldEdgeRelations().Create(ctx, edge)
	if err != nil {
		t.Fatal(err)
	}

	if edge.ID() == nil {
		t.Fatal("edge ID must not be nil after create")
	}

	assert.Check(t, edge.Meta.IsAdmin)
	assert.Check(t, edge.Meta.IsActive)

	results, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.
				FieldEdgeRelations(
					filter.EdgeRelation.CreatedAt.Before(time.Now().Add(time.Minute)),
				).
				SpecialTypes(
					filter.SpecialTypes.Name.Equal("edge_target"),
				),
		).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "edge_source", results[0].FieldString)
}
