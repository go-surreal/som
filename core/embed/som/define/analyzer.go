//go:build embed

package define

// Tokenizer represents a SurrealDB tokenizer type.
type Tokenizer string

const (
	// Blank tokenizer breaks down a text into tokens by creating a new token each
	// time it encounters a space, tab, or newline character. It’s a straightforward
	// way to split text into words or chunks based on whitespace.
	//
	// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#blank
	Blank Tokenizer = "blank"

	// Camel tokenizer is used for identifying and creating tokens when the next
	// character in the text is uppercase. This is particularly useful for processing
	// camelCase or PascalCase text, common in programming, to split them into meaningful words.
	//
	// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#camel
	Camel Tokenizer = "camel"

	// The class tokenizer segments text into tokens by detecting changes (digit, letter, punctuation, blank) in the Unicode class of characters. It creates a new token when the character class changes, distinguishing between digits, letters, punctuation, and blanks. This allows for flexible tokenization based on character types.
	//
	// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#class
	Class Tokenizer = "class"

	// The punct tokenizer generates tokens by breaking the text whenever a punctuation character is encountered. It’s suitable for tokenizing sentences or breaking text into smaller units based on punctuation marks.
	//
	// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#punct
	Punct Tokenizer = "punct"
)

// Filter represents a SurrealDB filter configuration.
type Filter struct {
	Name   string
	Params []any
}

var (
	// The ascii filter is responsible for processing tokens by replacing or removing diacritical marks (accents and special characters) from the text. It helps standardize text by converting accented characters to their basic ASCII equivalents, making it more suitable for various text analysis tasks.
	//
	// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#ascii
	Ascii = Filter{Name: "ascii"}

	// The lowercase filter converts tokens to lowercase, ensuring that text is consistently in lowercase format. This is often used to make text case-insensitive for search and analysis purposes.
	//
	// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#lowercase
	Lowercase = Filter{Name: "lowercase"}

	// The uppercase filter converts tokens to uppercase, ensuring text consistency in uppercase format. It can be useful when case-insensitivity is required for specific analysis or search operations.
	//
	// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#uppercase
	Uppercase = Filter{Name: "uppercase"}
)

// The edgengram filter is used to create tokens that represent prefixes of terms. It generates a sequence of tokens that gradually build up a term, which can be useful for autocomplete or searching based on partial words. It accepts two parameters min and max which define the minimum and maximum amount of characters in the prefix.
//
// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#edgengramminmax
func Edgengram(min, max int) Filter {
	return Filter{Name: "edgengram", Params: []any{min, max}}
}

// The mapping filter is designed to enable lemmatization within SurrealDB.
//
// Lemmatization is the process of reducing words to their base or dictionary form. The mapper mechanism allows users to specify a custom dictionary file that maps terms to their base forms. This dictionary file is then used by SurrealDB’s analyzer to standardize terms as they are indexed, improving search consistency.
//
// This is particularly useful for handling irregular verbs and other terms that the default “snowball” filter cannot handle. Lemmatization files are easy to put together and to find online, making it possible to customize full-text search for smaller languages.
//
// How does the mapper work?
//
// Configuration: In the SQL statement below, the mapper parameter is specified within the analyzer definition. This parameter points to the file that contains the term mappings for lemmatization.
//
// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#mapperpath
func Mapper(path string) Filter {
	return Filter{Name: "mapper", Params: []any{path}}
}

// The ngram filter is used to create a sequence of ‘n’ tokens from a given sample of text or speech. These items can be syllables, letters, words or base pairs according to the application. It accepts two parameters min and max which indicates that you want to create n-grams starting from min to size of max.
//
// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#ngramminmax
func Ngram(min, max int) Filter {
	return Filter{Name: "ngram", Params: []any{min, max}}
}

// The snowball filter applies Snowball stemming to tokens, reducing them to their root form and converts the case to lowercase.
//
// Docs: https://surrealdb.com/docs/surrealql/statements/define/analyzer#snowballlanguage
func Snowball(lang Language) Filter {
	return Filter{Name: "snowball", Params: []any{string(lang)}}
}

// Language represents a language for snowball stemming.
type Language string

const (
	Arabic     Language = "arabic"
	Danish     Language = "danish"
	Dutch      Language = "dutch"
	English    Language = "english"
	French     Language = "french"
	German     Language = "german"
	Greek      Language = "greek"
	Hungarian  Language = "hungarian"
	Italian    Language = "italian"
	Norwegian  Language = "norwegian"
	Portuguese Language = "portuguese"
	Romanian   Language = "romanian"
	Russian    Language = "russian"
	Spanish    Language = "spanish"
	Swedish    Language = "swedish"
	Tamil      Language = "tamil"
	Turkish    Language = "turkish"
)

// FulltextAnalyzerBuilder builds a DEFINE ANALYZER statement.
type FulltextAnalyzerBuilder struct {
	name       string
	tokenizers []Tokenizer
	filters    []Filter
	function   string
	comment    string
}

// FulltextAnalyzer creates a new analyzer definition.
func FulltextAnalyzer(name string) *FulltextAnalyzerBuilder {
	return &FulltextAnalyzerBuilder{name: name}
}

// Tokenizers sets the tokenizers for the analyzer.
func (b *FulltextAnalyzerBuilder) Tokenizers(tokenizers ...Tokenizer) *FulltextAnalyzerBuilder {
	b.tokenizers = tokenizers
	return b
}

// Filters sets the filters for the analyzer.
func (b *FulltextAnalyzerBuilder) Filters(filters ...Filter) *FulltextAnalyzerBuilder {
	b.filters = filters
	return b
}

// Function sets a custom function for the analyzer.
func (b *FulltextAnalyzerBuilder) Function(fn string) *FulltextAnalyzerBuilder {
	b.function = fn
	return b
}

// Comment adds a comment to the analyzer definition.
func (b *FulltextAnalyzerBuilder) Comment(comment string) *FulltextAnalyzerBuilder {
	b.comment = comment
	return b
}

// Getters for parser access
func (b *FulltextAnalyzerBuilder) GetName() string            { return b.name }
func (b *FulltextAnalyzerBuilder) GetTokenizers() []Tokenizer { return b.tokenizers }
func (b *FulltextAnalyzerBuilder) GetFilters() []Filter       { return b.filters }
func (b *FulltextAnalyzerBuilder) GetFunction() string        { return b.function }
func (b *FulltextAnalyzerBuilder) GetComment() string         { return b.comment }

// analyzerJSON is the JSON representation of an analyzer definition.
type analyzerJSON struct {
	Name       string       `json:"name"`
	Tokenizers []string     `json:"tokenizers"`
	Filters    []filterJSON `json:"filters"`
}

// filterJSON is the JSON representation of a filter.
type filterJSON struct {
	Name   string `json:"name"`
	Params []any  `json:"params,omitempty"`
}

// toJSON converts the analyzer builder to its JSON representation.
func (b *FulltextAnalyzerBuilder) toJSON() analyzerJSON {
	tokenizers := make([]string, len(b.tokenizers))
	for i, t := range b.tokenizers {
		tokenizers[i] = string(t)
	}

	filters := make([]filterJSON, len(b.filters))
	for i, f := range b.filters {
		filters[i] = filterJSON{
			Name:   f.Name,
			Params: f.Params,
		}
	}

	return analyzerJSON{
		Name:       b.name,
		Tokenizers: tokenizers,
		Filters:    filters,
	}
}
