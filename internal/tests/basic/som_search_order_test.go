package basic

import (
	"strings"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/gen/som/repo"
	"github.com/go-surreal/som/tests/basic/gen/som/where"
	"gotest.tools/v3/assert"
)

func TestFulltextSearchOrder(t *testing.T) {
	client := &repo.ClientImpl{}

	// Test 1: Score sort first, then field sort
	query1 := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		Order(query.Score(0).Desc(), by.AllFieldTypes.String.Asc())

	assert.Assert(t, strings.Contains(query1.Describe(),
		"ORDER BY __som_search_score_combined DESC, __som_sort__string ASC"))

	// Test 2: Field sort first, then score sort
	query2 := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		Order(by.AllFieldTypes.String.Asc(), query.Score(0).Desc())

	assert.Assert(t, strings.Contains(query2.Describe(),
		"ORDER BY __som_sort__string ASC, __som_search_score_combined DESC"))
}

func TestFulltextSearchScoreCombination(t *testing.T) {
	client := &repo.ClientImpl{}

	// Test Sum (default)
	q1 := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		Order(query.Score(0, 1).Desc())
	assert.Assert(t, strings.Contains(q1.Describe(),
		"(search::score(0) + search::score(1)) AS __som_search_score_combined"))

	// Test Sum (explicit)
	q1b := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		Order(query.Score(0, 1).Sum().Desc())
	assert.Assert(t, strings.Contains(q1b.Describe(),
		"(search::score(0) + search::score(1)) AS __som_search_score_combined"))

	// Test Max
	q2 := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		Order(query.Score(0, 1).Max().Desc())
	assert.Assert(t, strings.Contains(q2.Describe(),
		"math::max(search::score(0), search::score(1)) AS __som_search_score_combined"))

	// Test Average
	q3 := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		Order(query.Score(0, 1).Average().Desc())
	assert.Assert(t, strings.Contains(q3.Describe(),
		"((search::score(0) + search::score(1)) / 2) AS __som_search_score_combined"))

	// Test Weighted
	q4 := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test")).
		Order(query.Score(0, 1).Weighted(2.0, 0.5).Desc())
	assert.Assert(t, strings.Contains(q4.Describe(),
		"(search::score(0) * 2 + search::score(1) * 0.5) AS __som_search_score_combined"))
}

func TestSearchWithOffsets(t *testing.T) {
	client := &repo.ClientImpl{}

	q := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test").WithOffsets())

	assert.Assert(t, strings.Contains(q.Describe(),
		"search::offsets(0) AS __som_search_off_0"))
}

func TestSearchWithHighlightsAndOffsets(t *testing.T) {
	client := &repo.ClientImpl{}

	q := client.AllFieldTypesRepo().Query().
		Search(where.AllFieldTypes.String.Matches("test").
			WithHighlights("<b>", "</b>").
			WithOffsets())

	desc := q.Describe()
	assert.Assert(t, strings.Contains(desc, "search::highlight"))
	assert.Assert(t, strings.Contains(desc, "search::offsets"))
}

