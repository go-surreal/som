//go:build som

package model

import "github.com/go-surreal/som/define"

// Analyzer definitions for fulltext search

var english = define.Analyzer("english").
	Tokenizers(define.Blank, define.Punct).
	Filters(define.Lowercase, define.Snowball("en"))

var autocomplete = define.Analyzer("autocomplete").
	Tokenizers(define.Class).
	Filters(define.Lowercase, define.Edgengram(1, 10))

// Search configurations

var _ = define.Search("english_search").
	Analyzer(english).
	BM25(1.2, 0.75).
	Highlights()

var _ = define.Search("autocomplete_search").
	Analyzer(autocomplete)
