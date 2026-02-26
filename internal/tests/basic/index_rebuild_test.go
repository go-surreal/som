package basic

import (
	"context"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/model"
)

func TestIndexRebuildCount(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.AllTypesRepo().Index().Count().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexRebuildSearch(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString:   "hello world",
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.AllTypesRepo().Index().SearchFieldString().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexRebuildRegularIndex(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
		FieldCredentials: model.Credentials{
			Username: "testuser",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.AllTypesRepo().Index().IndexFieldCredentialsUsername().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexRebuildSoftDeleteIndex(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.SpecialTypesRepo().Index().IndexDeletedAt().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexRebuildMultiple(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString:   "test data",
			FieldTime:     time.Now(),
			FieldDuration: time.Second,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	idx := client.AllTypesRepo().Index()

	err := idx.Count().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = idx.SearchFieldString().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = idx.IndexFieldCredentialsUsername().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexRebuildMinimalNode(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.WeatherRepo().Index().Count().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
