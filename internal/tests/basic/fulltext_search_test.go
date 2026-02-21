package basic

import (
	"math"
	"strings"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/gen/som/repo"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestFullTextSearchBasic(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "the quick brown fox jumps over the lazy dog",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("quick fox")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "the quick brown fox jumps over the lazy dog", results[0].Model.FieldString)
}

func TestFullTextSearchNoMatch(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "hello world",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("nonexistent terms")).
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
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString: s,
		})
		if err != nil {
			t.Fatalf("failed to create test data: %v", err)
		}
	}

	results, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("programming")).
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

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "searchable content here",
		FieldInt:    42,
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "searchable content there",
		FieldInt:    100,
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("searchable")).
		Where(filter.AllTypes.FieldInt.Equal(42)).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
	assert.Equal(t, 42, results[0].Model.FieldInt)
}

func TestFullTextSearchQueryDescribe(t *testing.T) {
	client, cleanup := prepareDatabase(t.Context(), t)
	defer cleanup()

	query := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test query"))

	desc := query.Describe()
	t.Logf("Query: %s", desc)

	assert.Assert(t, len(desc) > 0)
}

func TestFullTextSearchWithRef(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "testing explicit ref",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("explicit").Ref(5)).
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

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "highlight this word please",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("highlight").WithHighlights("<mark>", "</mark>")).
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

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "apple pie is delicious",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "orange juice is refreshing",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "banana bread is tasty",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	// Search() now uses OR by default - matches documents with "apple" OR "orange"
	results, err := client.AllTypesRepo().Query().
		Search(
			filter.AllTypes.FieldString.Matches("apple"),
			filter.AllTypes.FieldString.Matches("orange"),
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

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "apple pie is delicious and sweet",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "apple juice is refreshing",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "orange juice is also refreshing",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	// SearchAll() uses AND - only matches documents with BOTH "apple" AND "delicious"
	results, err := client.AllTypesRepo().Query().
		SearchAll(
			filter.AllTypes.FieldString.Matches("apple"),
			filter.AllTypes.FieldString.Matches("delicious"),
		).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	// Only the first document contains both "apple" and "delicious"
	assert.Equal(t, 1, len(results))
	assert.Assert(t, results[0].Model.FieldString == "apple pie is delicious and sweet")
}

func TestFullTextSearchFirstMatch(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "first result here",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "second result here",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	result, ok, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("result")).
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

	_, ok, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("nonexistent")).
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

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "get all without metadata",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	models, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("metadata")).
		All(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(models))
	assert.Equal(t, "get all without metadata", models[0].FieldString)
}

func TestFullTextSearchScore(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "test test test repeated words",
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	results, err := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		AllMatches(ctx)

	if err != nil {
		t.Fatalf("failed to execute search: %v", err)
	}

	assert.Equal(t, 1, len(results))
	t.Logf("Search score: %f", results[0].Score())
	assert.Assert(t, len(results[0].Scores) > 0, "expected score projection to be present")
	assert.Assert(t, !math.IsNaN(results[0].Score()) && !math.IsInf(results[0].Score(), 0),
		"score must be a finite number, got %f", results[0].Score())
}

func TestFullTextSearchMultipleScoreSorts(t *testing.T) {
	ctx := t.Context()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	strVal := "test data for multiple scores"
	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString:    strVal,
		FieldStringPtr: &strVal,
	})
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	// Search on two fields (String and StringPtr) to get two search refs (0 and 1)
	// Then sort by Score(0) and Score(1) to get two different score aliases
	results, err := client.AllTypesRepo().Query().
		Search(
			filter.AllTypes.FieldString.Matches("test"),
			filter.AllTypes.FieldStringPtr.Matches("test"),
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

func TestFulltextSearchOrder(t *testing.T) {
	client := &repo.ClientImpl{}

	// Test 1: Score sort first, then field sort
	query1 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(query.Score(0).Desc(), by.AllTypes.FieldString.Asc())

	assert.Assert(t, strings.Contains(query1.Describe(),
		"ORDER BY __som__search_score_0 DESC, field_string ASC"))

	// Test 2: Field sort first, then score sort
	query2 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(by.AllTypes.FieldString.Asc(), query.Score(0).Desc())

	assert.Assert(t, strings.Contains(query2.Describe(),
		"ORDER BY field_string ASC, __som__search_score_0 DESC"))
}

func TestFulltextSearchScoreCombination(t *testing.T) {
	client := &repo.ClientImpl{}

	// Test Sum (default)
	q1 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(query.Score(0, 1).Desc())
	assert.Assert(t, strings.Contains(q1.Describe(),
		"(search::score(0) + search::score(1)) AS __som__search_score_0_1"))

	// Test Sum (explicit)
	q1b := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(query.Score(0, 1).Sum().Desc())
	assert.Assert(t, strings.Contains(q1b.Describe(),
		"(search::score(0) + search::score(1)) AS __som__search_score_0_1"))

	// Test Max
	q2 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(query.Score(0, 1).Max().Desc())
	assert.Assert(t, strings.Contains(q2.Describe(),
		"math::max(search::score(0), search::score(1)) AS __som__search_score_0_1"))

	// Test Average
	q3 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(query.Score(0, 1).Average().Desc())
	assert.Assert(t, strings.Contains(q3.Describe(),
		"((search::score(0) + search::score(1)) / 2) AS __som__search_score_0_1"))

	// Test Weighted
	q4 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(query.Score(0, 1).Weighted(2.0, 0.5).Desc())
	assert.Assert(t, strings.Contains(q4.Describe(),
		"(search::score(0) * 2 + search::score(1) * 0.5) AS __som__search_score_0_1"))
}

func TestSearchWithOffsets(t *testing.T) {
	client := &repo.ClientImpl{}

	q := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test").WithOffsets())

	assert.Assert(t, strings.Contains(q.Describe(),
		"search::offsets(0) AS __som__search_offsets_0"))
}

func TestSearchWithHighlightsAndOffsets(t *testing.T) {
	client := &repo.ClientImpl{}

	q := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test").
			WithHighlights("<b>", "</b>").
			WithOffsets())

	desc := q.Describe()
	assert.Assert(t, strings.Contains(desc, "search::highlight"))
	assert.Assert(t, strings.Contains(desc, "search::offsets"))
}

func TestFulltextSearchValidTypes(t *testing.T) {
	client := &repo.ClientImpl{}

	// Test string
	q1 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test"))
	assert.Assert(t, strings.Contains(q1.Describe(), "field_string @0@"))

	// Test *string
	q2 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldStringPtr.Matches("test"))
	assert.Assert(t, strings.Contains(q2.Describe(), "field_string_ptr @0@"))

	// Test []string (named "Other" in the model)
	q3 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldOther.Matches("test"))
	assert.Assert(t, strings.Contains(q3.Describe(), "field_other @0@"))

	// Test []*string
	q4 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldStringPtrSlice.Matches("test"))
	assert.Assert(t, strings.Contains(q4.Describe(), "field_string_ptr_slice @0@"))

	// Test *[]string
	q5 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldStringSlicePtr.Matches("test"))
	assert.Assert(t, strings.Contains(q5.Describe(), "field_string_slice_ptr @0@"))

	// Test *[]*string
	q6 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldStringPtrSlicePtr.Matches("test"))
	assert.Assert(t, strings.Contains(q6.Describe(), "field_string_ptr_slice_ptr @0@"))
}
