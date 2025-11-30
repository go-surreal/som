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
