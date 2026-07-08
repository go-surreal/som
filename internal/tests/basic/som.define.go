//go:build som

package basic

import (
	"som.test/gen/som/define"
	"som.test/gen/som/define/aggregate"
	"som.test/gen/som/filter"
	"som.test/model"
)

func Definitions() define.Definitions {
	return define.Definitions{
		Searches: []*define.SearchBuilder{
			searchEnglish,
			searchAutocomplete,
		},
		Views: []define.ViewDefinition{
			allTypesSummary,
		},
	}
}

//
// -- ANALYZERS
//

var (
	english = define.FulltextAnalyzer("english").
		Tokenizers(define.Blank, define.Punct).
		Filters(define.Lowercase, define.Snowball(define.English))

	autocomplete = define.FulltextAnalyzer("autocomplete").
		Tokenizers(define.Class).
		Filters(define.Lowercase, define.Edgengram(1, 10))
)

//
// -- SEARCHES
//

var (
	searchEnglish = define.Search("english_search").
		FulltextAnalyzer(english).
		BM25(1.2, 0.75).
		Highlights()

	searchAutocomplete = define.Search("autocomplete_search").
		FulltextAnalyzer(autocomplete)
)

//
// -- VIEWS
//

var allTypesSummary = define.View[model.AllTypesSummary, model.AllTypes]().
	Project(
		define.As(filter.AllTypesSummary.Category, filter.AllTypes.FieldString),
		define.As(filter.AllTypesSummary.Total, aggregate.Count(filter.AllTypes.FieldString)),
		define.As(filter.AllTypesSummary.AvgValue, aggregate.Mean(filter.AllTypes.FieldFloat64)),
	).
	GroupBy(filter.AllTypes.FieldString)
