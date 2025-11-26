//go:build embed

// Package define provides types for defining SurrealDB analyzers and search
// configurations. These types are used in files with the "//go:build som"
// build tag and are parsed by the som code generator at generation time.
// They have no runtime effect.
package define

//
// -- ANALYZER BUILDER
//

// AnalyzerBuilder is used to define a custom analyzer for fulltext search.
// The analyzer name and configuration are extracted by the som code generator.
type AnalyzerBuilder struct {
	name       string
	tokenizers []Tokenizer
	filters    []Filter
}

// Analyzer creates a new analyzer definition with the given name.
// The name is used to reference this analyzer in Search configurations.
//
// Example:
//
//	var english = define.Analyzer("english").
//	    Tokenizers(define.Blank, define.Punct).
//	    Filters(define.Lowercase, define.Snowball("en"))
func Analyzer(name string) *AnalyzerBuilder {
	return &AnalyzerBuilder{name: name}
}

// Tokenizers sets the tokenizers for this analyzer.
// Tokenizers split input text into tokens for indexing.
func (a *AnalyzerBuilder) Tokenizers(tokenizers ...Tokenizer) *AnalyzerBuilder {
	a.tokenizers = tokenizers
	return a
}

// Filters sets the filters for this analyzer.
// Filters transform tokens after tokenization (e.g., lowercase, stemming).
func (a *AnalyzerBuilder) Filters(filters ...Filter) *AnalyzerBuilder {
	a.filters = filters
	return a
}

// Name returns the analyzer's name (used by parser).
func (a *AnalyzerBuilder) Name() string {
	return a.name
}

//
// -- SEARCH BUILDER
//

// SearchBuilder is used to define a search configuration that combines
// an analyzer with search options like BM25 parameters and highlighting.
type SearchBuilder struct {
	name       string
	analyzer   *AnalyzerBuilder
	bm25K1     float64
	bm25B      float64
	bm25Set    bool
	highlights bool
}

// Search creates a new search configuration with the given name.
// The name is used to reference this configuration in struct tags.
//
// Example:
//
//	var englishSearch = define.Search("english_search").
//	    Analyzer(english).
//	    BM25(1.2, 0.75).
//	    Highlights()
//
// Then use in struct tags:
//
//	Title string `som:"search:english_search"`
func Search(name string) *SearchBuilder {
	return &SearchBuilder{name: name}
}

// Analyzer sets the analyzer to use for this search configuration.
// Pass the variable from an Analyzer() definition.
func (s *SearchBuilder) Analyzer(a *AnalyzerBuilder) *SearchBuilder {
	s.analyzer = a
	return s
}

// BM25 sets the BM25 ranking parameters.
// - k1: Term frequency saturation parameter (typically 1.2)
// - b: Document length normalization parameter (typically 0.75)
func (s *SearchBuilder) BM25(k1, b float64) *SearchBuilder {
	s.bm25K1 = k1
	s.bm25B = b
	s.bm25Set = true
	return s
}

// Highlights enables search result highlighting for this configuration.
func (s *SearchBuilder) Highlights() *SearchBuilder {
	s.highlights = true
	return s
}

// Name returns the search config name (used by parser).
func (s *SearchBuilder) Name() string {
	return s.name
}

// AnalyzerRef returns the referenced analyzer (used by parser).
func (s *SearchBuilder) AnalyzerRef() *AnalyzerBuilder {
	return s.analyzer
}

// HasBM25 returns whether BM25 parameters are set (used by parser).
func (s *SearchBuilder) HasBM25() bool {
	return s.bm25Set
}

// BM25Params returns the BM25 k1 and b parameters (used by parser).
func (s *SearchBuilder) BM25Params() (k1, b float64) {
	return s.bm25K1, s.bm25B
}

// HasHighlights returns whether highlighting is enabled (used by parser).
func (s *SearchBuilder) HasHighlights() bool {
	return s.highlights
}

//
// -- TOKENIZERS
//

// Tokenizer represents a tokenization strategy for text analysis.
type Tokenizer struct {
	name string
}

// Name returns the tokenizer's name (used by parser).
func (t Tokenizer) Name() string {
	return t.name
}

// Available tokenizers for SurrealDB analyzers.
var (
	// Blank tokenizes on whitespace characters.
	Blank = Tokenizer{name: "blank"}

	// Camel tokenizes on camelCase boundaries.
	Camel = Tokenizer{name: "camel"}

	// Class tokenizes based on Unicode character classes.
	Class = Tokenizer{name: "class"}

	// Punct tokenizes on punctuation characters.
	Punct = Tokenizer{name: "punct"}
)

//
// -- FILTERS
//

// Filter represents a token filter for text analysis.
type Filter struct {
	name   string
	params []any
}

// Name returns the filter's name (used by parser).
func (f Filter) Name() string {
	return f.name
}

// Params returns the filter's parameters (used by parser).
func (f Filter) Params() []any {
	return f.params
}

// Available filters for SurrealDB analyzers.
var (
	// Ascii normalizes tokens to ASCII characters.
	Ascii = Filter{name: "ascii"}

	// Lowercase converts tokens to lowercase.
	Lowercase = Filter{name: "lowercase"}

	// Uppercase converts tokens to uppercase.
	Uppercase = Filter{name: "uppercase"}
)

// Snowball creates a stemming filter for the specified language.
// Supported languages: "en" (English), "de" (German), etc.
func Snowball(lang string) Filter {
	return Filter{name: "snowball", params: []any{lang}}
}

// Edgengram creates an edge n-gram filter with min and max lengths.
// This is useful for autocomplete/typeahead functionality.
func Edgengram(min, max int) Filter {
	return Filter{name: "edgengram", params: []any{min, max}}
}

// Ngram creates an n-gram filter with min and max lengths.
func Ngram(min, max int) Filter {
	return Filter{name: "ngram", params: []any{min, max}}
}
