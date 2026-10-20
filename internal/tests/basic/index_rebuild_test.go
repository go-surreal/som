package basic

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/surrealdb/surrealdb.go"
	"som.test/gen/som"
	"som.test/gen/som/index"
	"som.test/gen/som/repo"
	"som.test/model"
)

const (
	indexReadyTimeout      = 5 * time.Second
	indexReadyPollInterval = 20 * time.Millisecond
)

// waitForIndexReady blocks until the concurrent build of the given index
// has finished, since SurrealDB refuses to rebuild an index while it is
// still building.
func waitForIndexReady(ctx context.Context, t *testing.T, client repo.Client, table, index string) {
	t.Helper()

	query := fmt.Sprintf("RETURN [(INFO FOR INDEX %s ON TABLE %s).building.status]", index, table)
	deadline := time.Now().Add(indexReadyTimeout)

	for time.Now().Before(deadline) {
		result, err := client.Raw(ctx, query, nil)
		if err != nil {
			t.Fatalf("could not query build status of index %s: %v", index, err)
		}

		statuses, err := result.Scan[string]()
		if err != nil {
			t.Fatalf("could not decode build status of index %s: %v", index, err)
		}

		if len(statuses) == 1 && statuses[0] == "ready" {
			return
		}

		time.Sleep(indexReadyPollInterval)
	}

	t.Fatalf("index %s was not ready within %s", index, indexReadyTimeout)
}

func TestIndexRebuildCount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
		FieldMonth:    time.January,
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
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString:   "hello world",
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
		FieldMonth:    time.January,
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
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
		FieldMonth:    time.January,
		FieldCredentials: model.Credentials{
			Username: "testuser",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	waitForIndexReady(ctx, t, client, "all_types", "__som__all_types_index_field_credentials_username")

	err = client.AllTypesRepo().Index().IndexFieldCredentialsUsername().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexRebuildSoftDeleteIndex(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	waitForIndexReady(ctx, t, client, "special_types", "__som__special_types_deleted_at")

	err := client.SpecialTypesRepo().Index().DeletedAt().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

// buildingDatabase answers every query the way SurrealDB answers a rebuild
// of an index whose build is still running.
type buildingDatabase struct{}

func (buildingDatabase) Query(context.Context, string, map[string]any) ([]byte, error) {
	queryErr := &surrealdb.QueryError{
		Message: "Database index `__som__special_types_deleted_at` is currently building",
	}

	return nil, fmt.Errorf("query statement 0 failed: %w", queryErr)
}

func TestIndexRebuildWhileBuilding(t *testing.T) {
	t.Parallel()

	err := index.NewSpecialTypes(buildingDatabase{}).DeletedAt().Rebuild(context.Background())

	if !errors.Is(err, som.ErrIndexBuilding) {
		t.Fatalf("expected som.ErrIndexBuilding, got: %v", err)
	}

	var queryErr *surrealdb.QueryError
	if !errors.As(err, &queryErr) {
		t.Fatalf("expected the query error to stay reachable, got: %v", err)
	}
}

func TestIndexRebuildMultiple(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString:   "test data",
			FieldTime:     time.Now(),
			FieldDuration: time.Second,
			FieldMonth:    time.January,
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

	waitForIndexReady(ctx, t, client, "all_types", "__som__all_types_index_field_credentials_username")

	err = idx.IndexFieldCredentialsUsername().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexRebuildMinimalNode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.WeatherRepo().Index().Count().Rebuild(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
