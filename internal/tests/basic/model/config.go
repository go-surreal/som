//go:build som

package model

import "github.com/go-surreal/som/tests/basic/gen/som/define"

// Analyzer definitions for fulltext search

var english = define.FulltextAnalyzer("english").
	Tokenizers(define.Blank, define.Punct).
	Filters(define.Lowercase, define.Snowball(define.English))

var autocomplete = define.FulltextAnalyzer("autocomplete").
	Tokenizers(define.Class).
	Filters(define.Lowercase, define.Edgengram(1, 10))

// Search configurations

var _ = define.Search("english_search").
	FulltextAnalyzer(english).
	BM25(1.2, 0.75).
	Highlights()

var _ = define.Search("autocomplete_search").
	FulltextAnalyzer(autocomplete)

type X struct {
	som.Model

	Field1 string `som:"index"`
}
