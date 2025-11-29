//go:build som

package model

import "github.com/go-surreal/som/tests/basic/gen/som/define"

func Definitions() define.Definitions {
	return define.Definitions{
		Searches: []*define.SearchBuilder{
			searchEnglish,
			searchAutocomplete,
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
