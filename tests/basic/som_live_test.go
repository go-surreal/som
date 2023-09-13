package basic

import (
	"context"
	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/gen/som/where"

	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"testing"
	"time"
)

func TestCreateWithFieldsLikeDBResponse(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	newModel := &model.FieldsLikeDBResponse{
		Status: "some value",
	}

	err := client.FieldsLikeDBResponseRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	readModel, exists, err := client.FieldsLikeDBResponseRepo().Read(ctx, newModel.ID())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, exists)
	assert.Equal(t, "some value", readModel.Status)

	readModel.Status = "some other value"

	err = client.FieldsLikeDBResponseRepo().Update(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "some other value", readModel.Status)

	err = client.FieldsLikeDBResponseRepo().Delete(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLiveQueries(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	newModel := &model.FieldsLikeDBResponse{
		Status: "some value",
	}

	liveChan, err := client.FieldsLikeDBResponseRepo().Query().Live(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = client.FieldsLikeDBResponseRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	newModel.Status = "some other value"
	err = client.FieldsLikeDBResponseRepo().Update(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	err = client.FieldsLikeDBResponseRepo().Delete(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	// LIVE CREATE

	liveRes, more := <-liveChan
	if !more {
		t.Fatal("liveChan closed unexpectedly")
	}

	liveCreate, ok := liveRes.(query.LiveCreate[*model.FieldsLikeDBResponse])
	if !ok {
		t.Fatal("liveChan did not receive a create event")
	}

	created, err := liveCreate.Get()
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, is.Equal(newModel.ID(), created.ID()))

	// LIVE UPDATE

	liveRes, more = <-liveChan
	if !more {
		t.Fatal("liveChan closed unexpectedly")
	}

	liveUpdate, ok := liveRes.(query.LiveUpdate[*model.FieldsLikeDBResponse])
	if !ok {
		t.Fatal("liveChan did not receive an update event")
	}

	updated, err := liveUpdate.Get()
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, is.Equal(newModel.ID(), updated.ID()))

	// LIVE DELETE

	liveRes, more = <-liveChan
	if !more {
		t.Fatal("liveChan closed unexpectedly")
	}

	liveDelete, ok := liveRes.(query.LiveDelete)
	if !ok {
		t.Fatal("liveChan did not receive a delete event")
	}

	deletedID, err := liveDelete.Get()
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, is.Equal("fields_like_db_response:"+newModel.ID(), deletedID))
}

func TestLiveQueriesFilter(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	liveChan, err := client.FieldsLikeDBResponseRepo().Query().
		Filter(
			where.FieldsLikeDBResponse.Status.In([]string{"some value", "some other value"}),
		).
		Live(ctx)
	if err != nil {
		t.Fatal(err)
	}

	newModel1 := &model.FieldsLikeDBResponse{
		Status: "some value",
	}

	err = client.FieldsLikeDBResponseRepo().Create(ctx, newModel1)
	if err != nil {
		t.Fatal(err)
	}

	newModel2 := &model.FieldsLikeDBResponse{
		Status: "some unsupported value",
	}

	err = client.FieldsLikeDBResponseRepo().Create(ctx, newModel2)
	if err != nil {
		t.Fatal(err)
	}

	newModel3 := &model.FieldsLikeDBResponse{
		Status: "some other value",
	}

	err = client.FieldsLikeDBResponseRepo().Create(ctx, newModel3)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"some value",
		// note: "some unsupported value" should not be received ;)
		"some other value",
	}

	for _ = range expected {
		select {

		case _, more := <-liveChan:
			{
				if !more {
					t.Fatal("liveChan closed unexpectedly")
				}

				// liveCreate, ok := liveRes.(query.LiveCreate[*model.FieldsLikeDBResponse])
				// if !ok {
				// 	t.Fatal("liveChan did not receive a create event")
				// }
				//
				// created, err := liveCreate.Get()
				// if err != nil {
				// 	t.Fatal(err)
				// }
				//
				// assert.Check(t, is.Equal(status, created.Status))

				t.Fatal("for beta.12 live queries with filters should not work yet")
			}

		case <-time.After(1 * time.Second):
			// t.Fatal("timeout waiting for live event")
			t.Log("correct, becuase live queries with filters are not supported yet")
		}
	}
}
