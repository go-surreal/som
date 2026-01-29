package basic

import (
	"strings"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/query"
	"github.com/go-surreal/som/tests/basic/gen/som/repo"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"gotest.tools/v3/assert"
)

func TestFulltextSearchOrder(t *testing.T) {
	client := &repo.ClientImpl{}

	// Test 1: Score sort first, then field sort
	query1 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(query.Score(0).Desc(), by.AllTypes.FieldString.Asc())

	assert.Assert(t, strings.Contains(query1.Describe(),
		"ORDER BY __som__search_score_0 DESC, __som__sort_field_string ASC"))

	// Test 2: Field sort first, then score sort
	query2 := client.AllTypesRepo().Query().
		Search(filter.AllTypes.FieldString.Matches("test")).
		Order(by.AllTypes.FieldString.Asc(), query.Score(0).Desc())

	assert.Assert(t, strings.Contains(query2.Describe(),
		"ORDER BY __som__sort_field_string ASC, __som__search_score_0 DESC"))
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

