package basic

import (
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
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
		Search(filter.AllFieldTypes.String.Matches("quick fox")).
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
		Search(filter.AllFieldTypes.String.Matches("nonexistent terms")).
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
		Search(filter.AllFieldTypes.String.Matches("programming")).
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
		Search(filter.AllFieldTypes.String.Matches("searchable")).
		Where(filter.AllFieldTypes.Int.Equal(42)).
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
		Search(filter.AllFieldTypes.String.Matches("test query"))

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
		Search(filter.AllFieldTypes.String.Matches("explicit").Ref(5)).
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
		Search(filter.AllFieldTypes.String.Matches("highlight").WithHighlights("<mark>", "</mark>")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
}

// TestFullTextSearchOrDefault tests that Search() combines conditions with OR by default.
// This is the standard search engine behavior where any matching term is sufficient.
func TestFullTextSearchOrDefault(t *testing.T) {
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

	// Search() now uses OR by default - matches documents with "apple" OR "orange"
	results, err := client.AllFieldTypesRepo().Query().
		Search(
			filter.AllFieldTypes.String.Matches("apple"),
			filter.AllFieldTypes.String.Matches("orange"),
		).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 2, len(results))
}

// TestFullTextSearchAndExplicit tests that SearchAll() combines conditions with AND.
// Use this when documents must match ALL search terms.
func TestFullTextSearchAndExplicit(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "apple pie is delicious and sweet",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "apple juice is refreshing",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String: "orange juice is also refreshing",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	// SearchAll() uses AND - only matches documents with BOTH "apple" AND "delicious"
	results, err := client.AllFieldTypesRepo().Query().
		SearchAll(
			filter.AllFieldTypes.String.Matches("apple"),
			filter.AllFieldTypes.String.Matches("delicious"),
		).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	// Only the first document contains both "apple" and "delicious"
	assert.Equal(t, 1, len(results))
	assert.Assert(t, results[0].Model.String == "apple pie is delicious and sweet")
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
		Search(filter.AllFieldTypes.String.Matches("result")).
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
		Search(filter.AllFieldTypes.String.Matches("nonexistent")).
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
		Search(filter.AllFieldTypes.String.Matches("metadata")).
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
		Search(filter.AllFieldTypes.String.Matches("test")).
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

func TestFullTextSearchMultipleScoreSorts(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	strVal := "test data for multiple scores"
	err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{
		String:    strVal,
		StringPtr: &strVal,
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	// Search on two fields (String and StringPtr) to get two search refs (0 and 1)
	// Then sort by Score(0) and Score(1) to get two different score aliases
	results, err := client.AllFieldTypesRepo().Query().
		Search(
			filter.AllFieldTypes.String.Matches("test"),
			filter.AllFieldTypes.StringPtr.Matches("test"),
		).
		Order(query.Score(0).Desc(), query.Score(1).Asc()).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
	t.Logf("Scores: %v", results[0].Scores)
	// We should have at least 2 scores:
	// - search::score(0) from first search clause
	// - search::score(1) from second search clause
	// - Plus the score projections from Order (which may duplicate the above)
	assert.Assert(t, len(results[0].Scores) >= 2,
		"expected at least 2 scores, got %d", len(results[0].Scores))
}
