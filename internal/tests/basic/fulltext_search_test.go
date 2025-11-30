package basic

import (
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/where"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestFullTextSearchBasic(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "the quick brown fox jumps over the lazy dog",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("quick fox")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "the quick brown fox jumps over the lazy dog", results[0].Model.String)
}

func TestFullTextSearchNoMatch(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "hello world",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("nonexistent terms")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 0, len(results))
}

func TestFullTextSearchMultipleResults(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	testData := []string{
		"golang is a programming language",
		"programming in go is fun",
		"python is another programming language",
	}

	for _, s := range testData {
		err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
			String: s,
		})
		if err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}

	results, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("programming")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 3, len(results))
}

func TestFullTextSearchWithFilter(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "searchable content here",
		Int:    42,
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "searchable content there",
		Int:    100,
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("searchable")).
		Filter(where.AllFieldTypes.Int.Equal(42)).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
	assert.Equal(t, 42, results[0].Model.Int)
}

func TestFullTextSearchQueryDescribe(t *testing.T) {
	client, cleanup := prepareDatabase(t.Context(), t)
	defer cleanup()

	query := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test query"))

	desc := query.Describe()
	t.Logf("Query: %s", desc)

	assert.Assert(t, len(desc) > 0)
}

func TestFullTextSearchWithRef(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "testing explicit ref",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("explicit").Ref(5)).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
}

func TestFullTextSearchWithHighlights(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "highlight this word please",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("highlight").WithHighlights("<mark>", "</mark>")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
}

func TestFullTextSearchAny(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "apple pie is delicious",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "orange juice is refreshing",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "banana bread is tasty",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllFieldTypesRepo().Query().
		SearchAny(
			where.AllFieldTypes.String.Matches("apple"),
			where.AllFieldTypes.String.Matches("orange"),
		).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 2, len(results))
}

func TestFullTextSearchFirstMatch(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "first result here",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "second result here",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	result, ok, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("result")).
		FirstMatch(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Assert(t, ok)
	assert.Assert(t, result.Model != nil)
}

func TestFullTextSearchFirstMatchNoResult(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	_, ok, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("nonexistent")).
		FirstMatch(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Assert(t, !ok)
}

func TestFullTextSearchAll(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "get all without metadata",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	models, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("metadata")).
		All(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(models))
	assert.Equal(t, "get all without metadata", models[0].String)
}

func TestFullTextSearchScore(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "test test test repeated words",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
	score := results[0].Score()
	t.Logf("Search score: %f", score)
	// BM25 scores can be negative, just verify we got a score value
	assert.Assert(t, score != 0, "expected non-zero search score")
}
