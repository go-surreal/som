package basic

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/gen/som/with"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestCreateWithAllTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	newModel := &model.AllTypes{
		FieldHookStatus: "some value",
		FieldMonth:      time.January,
	}

	err := client.AllTypesRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	readModel, exists, err := client.AllTypesRepo().Read(ctx, string(newModel.ID()))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, exists)
	assert.Equal(t, "[created]some value", readModel.FieldHookStatus)

	readModel.FieldHookStatus = "some other value"

	err = client.AllTypesRepo().Update(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "[updated]some other value", readModel.FieldHookStatus)

	err = client.AllTypesRepo().Delete(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLiveQueries(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	newModel := &model.AllTypes{
		FieldHookStatus: "some value",
		FieldMonth:      time.January,
	}

	liveChan, err := client.AllTypesRepo().Query().Live(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = client.AllTypesRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	newModel.FieldHookStatus = "some other value"
	err = client.AllTypesRepo().Update(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	err = client.AllTypesRepo().Delete(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	// LIVE CREATE

	var liveRes query.LiveResult[*model.AllTypes]
	var more bool

	select {
	case liveRes, more = <-liveChan:
		if !more {
			t.Fatal("liveChan closed unexpectedly")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for CREATE event")
	}

	liveCreate, ok := liveRes.(query.LiveCreate[*model.AllTypes])
	if !ok {
		t.Fatalf("expected LiveCreate event, got %T", liveRes)
	}

	created, err := liveCreate.Get()
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, is.Equal(newModel.ID(), created.ID()))

	// LIVE UPDATE

	select {
	case liveRes, more = <-liveChan:
		if !more {
			t.Fatal("liveChan closed unexpectedly")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for UPDATE event")
	}

	liveUpdate, ok := liveRes.(query.LiveUpdate[*model.AllTypes])
	if !ok {
		t.Fatalf("expected LiveUpdate event, got %T", liveRes)
	}

	updated, err := liveUpdate.Get()
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, is.Equal(newModel.ID(), updated.ID()))

	// LIVE DELETE

	select {
	case liveRes, more = <-liveChan:
		if !more {
			t.Fatal("liveChan closed unexpectedly")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for DELETE event")
	}

	liveDelete, ok := liveRes.(query.LiveDelete[*model.AllTypes])
	if !ok {
		t.Fatalf("expected LiveDelete event, got %T", liveRes)
	}

	deleted, err := liveDelete.Get()
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, is.Equal(newModel.ID(), deleted.ID()))

	// Test the automatic closing of the live channel when the context is canceled:

	cancel()

	select {

	case _, more := <-liveChan:
		if more {
			t.Fatal("liveChan did not close after context was canceled")
		}

	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for live channel to close after context was canceled")
	}
}

func TestLiveQueriesFilter(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	liveChan, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldHookStatus.In([]string{"[created]some value", "[created]some other value"}),
		).
		Live(ctx)
	if err != nil {
		t.Fatal(err)
	}

	newModel1 := &model.AllTypes{
		FieldHookStatus: "some value",
		FieldMonth:      time.January,
	}

	err = client.AllTypesRepo().Create(ctx, newModel1)
	if err != nil {
		t.Fatal(err)
	}

	newModel2 := &model.AllTypes{
		FieldHookStatus: "some unsupported value",
		FieldMonth:      time.January,
	}

	err = client.AllTypesRepo().Create(ctx, newModel2)
	if err != nil {
		t.Fatal(err)
	}

	newModel3 := &model.AllTypes{
		FieldHookStatus: "some other value",
		FieldMonth:      time.January,
	}

	err = client.AllTypesRepo().Create(ctx, newModel3)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"[created]some value",
		//"some unsupported value", // this one should be filtered out
		"[created]some other value",
	}

	for _, status := range expected {
		select {

		case liveRes, more := <-liveChan:
			{
				if !more {
					t.Fatal("liveChan closed unexpectedly")
				}

				liveCreate, ok := liveRes.(query.LiveCreate[*model.AllTypes])
				if !ok {
					t.Fatal("liveChan did not receive a create event")
				}

				created, err := liveCreate.Get()
				if err != nil {
					t.Fatal(err)
				}

				assert.Check(t, is.Equal(status, created.FieldHookStatus))
			}

		case <-time.After(10 * time.Second):
			t.Fatal("timeout waiting for live event")
		}
	}
}

func TestLiveQueryCount(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	liveCount, err := client.AllTypesRepo().Query().LiveCount(ctx)
	if err != nil {
		t.Fatal(err)
	}

	count := rand.Intn(randMax-randMin) + randMin

	var models []*model.AllTypes

	for i := 0; i < count; i++ {
		newModel := &model.AllTypes{
			FieldTime:     time.Now(),
			FieldDuration: time.Second,
			FieldMonth:    time.January,
		}

		if err := client.AllTypesRepo().Create(ctx, newModel); err != nil {
			t.Fatal(err)
		}

		models = append(models, newModel)
	}

	for i := 0; i <= count; i++ {
		assert.Equal(t, i, <-liveCount)
	}

	for _, delModel := range models {
		if err := client.AllTypesRepo().Delete(ctx, delModel); err != nil {
			t.Fatal(err)
		}
	}

	for i := count; i > 0; i-- {
		assert.Equal(t, i-1, <-liveCount)
	}

	select {

	case <-liveCount:
		t.Fatal("liveCount should not receive any more messages")

	case <-time.After(1 * time.Second):
	}

	// Test the automatic closing of the live channel when the context is canceled:

	cancel()

	select {

	case _, more := <-liveCount:
		if more {
			t.Fatal("liveCount did not close after context was canceled")
		}

	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for live channel to close after context was canceled")
	}
}

func TestLiveQueryWithFetch(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a group first
	group := &model.SpecialTypes{
		Name: "Test Group",
	}
	err := client.SpecialTypesRepo().Create(ctx, group)
	if err != nil {
		t.Fatal(err)
	}

	// Start live query with Fetch
	liveChan, err := client.AllTypesRepo().Query().
		Fetch(with.AllTypes.FieldNode()).
		Live(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Create a record with MainGroup set
	newModel := &model.AllTypes{
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
		FieldMonth:    time.January,
		FieldNode:     *group,
	}

	err = client.AllTypesRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	// Receive the live create event
	select {
	case liveRes, more := <-liveChan:
		if !more {
			t.Fatal("liveChan closed unexpectedly")
		}

		liveCreate, ok := liveRes.(query.LiveCreate[*model.AllTypes])
		if !ok {
			t.Fatalf("expected LiveCreate event, got %T", liveRes)
		}

		created, err := liveCreate.Get()
		if err != nil {
			t.Fatal(err)
		}

		assert.Check(t, is.Equal(newModel.ID(), created.ID()))
		// Verify that the fetched MainGroup has data populated
		assert.Check(t, is.Equal(group.Name, created.FieldNode.Name))

	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for CREATE event")
	}
}
