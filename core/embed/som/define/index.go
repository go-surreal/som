//go:build embed

package define

// SearchBuilder builds a SEARCH ANALYZER index configuration.
type SearchBuilder struct {
	name         string
	analyzer     *FulltextAnalyzerBuilder
	bm25K1       float64
	bm25B        float64
	hasBM25      bool
	highlights   bool
	concurrently bool
}

// Search creates a new fulltext search index configuration.
func Search(name string) *SearchBuilder {
	return &SearchBuilder{name: name}
}

// FulltextAnalyzer sets the analyzer for the search index.
func (b *SearchBuilder) FulltextAnalyzer(analyzer *FulltextAnalyzerBuilder) *SearchBuilder {
	b.analyzer = analyzer
	return b
}

// BM25 sets BM25 ranking parameters.
func (b *SearchBuilder) BM25(k1, bParam float64) *SearchBuilder {
	b.bm25K1 = k1
	b.bm25B = bParam
	b.hasBM25 = true
	return b
}

// Highlights enables search result highlighting.
func (b *SearchBuilder) Highlights() *SearchBuilder {
	b.highlights = true
	return b
}

// Concurrently enables concurrent index building.
func (b *SearchBuilder) Concurrently() *SearchBuilder {
	b.concurrently = true
	return b
}

// Getters for parser access
func (b *SearchBuilder) GetName() string                           { return b.name }
func (b *SearchBuilder) GetFulltextAnalyzer() *FulltextAnalyzerBuilder { return b.analyzer }
func (b *SearchBuilder) GetBM25K1() float64                        { return b.bm25K1 }
func (b *SearchBuilder) GetBM25B() float64                         { return b.bm25B }
func (b *SearchBuilder) HasBM25() bool                             { return b.hasBM25 }
func (b *SearchBuilder) HasHighlights() bool                       { return b.highlights }
func (b *SearchBuilder) IsConcurrently() bool                      { return b.concurrently }
