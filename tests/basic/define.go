package basic

import (
	"github.com/go-surreal/som/core/embed/som"
	"github.com/go-surreal/som/tests/basic/gen/som/define"
)

var mainAnalyzer = som.Define().Analyzer("some_analyzer").
	Tokenizers(
		define.Blank(),
		define.Camel(),
		define.Class(),
		define.Punct(),
	).
	Filters(
		define.Ascii(),
		define.Edgengram(1, 3),
		define.Lowercase(),
		define.Uppercase(),
		define.Snowball(define.EN, define.DE),
	)

func init() {

	som.Define().Model().User().Index("search").
		On(
			field.User.Username,
			field.User.FirstName,
			field.User.LastName,
		).
		// Unique().
		Search(
			mainAnalyzer,
			define.BM25(k1, b),
			define.WithHighlights,
		)

}

// Legendary
// Exceptional
// Great
// Good
// Average
// Mediocre
// Poor
// Bad
// Terrible
// Abysmal
// Unplayable
