//go:build embed

package define

import "encoding/json"

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

// searchJSON is the JSON representation of a search configuration.
type searchJSON struct {
	Name         string  `json:"name"`
	AnalyzerName string  `json:"analyzer_name"`
	BM25K1       float64 `json:"bm25_k1,omitempty"`
	BM25B        float64 `json:"bm25_b,omitempty"`
	HasBM25      bool    `json:"has_bm25"`
	Highlights   bool    `json:"highlights"`
	Concurrently bool    `json:"concurrently"`
}

// toJSON converts the search builder to its JSON representation.
func (b *SearchBuilder) toJSON() searchJSON {
	analyzerName := ""
	if b.analyzer != nil {
		analyzerName = b.analyzer.name
	}
	return searchJSON{
		Name:         b.name,
		AnalyzerName: analyzerName,
		BM25K1:       b.bm25K1,
		BM25B:        b.bm25B,
		HasBM25:      b.hasBM25,
		Highlights:   b.highlights,
		Concurrently: b.concurrently,
	}
}

// Definitions holds all user-defined configurations.
type Definitions struct {
	Searches []*SearchBuilder
}

// defineOutputJSON is the JSON structure for all definitions.
type defineOutputJSON struct {
	Analyzers []analyzerJSON `json:"analyzers"`
	Searches  []searchJSON   `json:"searches"`
}

// ToJSON serializes all definitions to JSON.
func (d Definitions) ToJSON() ([]byte, error) {
	output := defineOutputJSON{}

	// Collect unique analyzers from searches
	seen := make(map[string]bool)
	for _, s := range d.Searches {
		if s.analyzer != nil && !seen[s.analyzer.name] {
			seen[s.analyzer.name] = true
			output.Analyzers = append(output.Analyzers, s.analyzer.toJSON())
		}
		output.Searches = append(output.Searches, s.toJSON())
	}

	return json.Marshal(output)
}
