//go:build embed

package lib

import (
	"strconv"
	"strings"
)

// Search represents a full-text search condition.
// This is a separate type from Filter to ensure type safety -
// Search conditions can only be used with Search()/SearchAll() methods.
type Search[M any] interface {
	Ref(ref int) Search[M]
	WithHighlights(prefix, suffix string) Search[M]
	WithOffsets() Search[M]
	build(ctx *context, ref int) SearchClause
}

// SearchClause holds the rendered search condition and metadata.
type SearchClause struct {
	SQL        string
	Ref        int
	Highlights bool
	HLPrefix   string
	HLSuffix   string
	Offsets    bool
}

type search[M any] struct {
	key        Key[M]
	terms      string
	ref        *int
	hlPrefix   string
	hlSuffix   string
	highlights bool
	offsets    bool
	boolOr     bool
}

// NewSearch creates a full-text search condition for the given key and terms.
// This is called by generated code for search-indexed string fields.
func NewSearch[M any](key Key[M], terms string) Search[M] {
	return &search[M]{key: key, terms: terms}
}

// NewSearchOr creates a full-text search condition that uses OR boolean mode.
// With OR mode, documents matching ANY of the search terms will be returned,
// rather than requiring ALL terms to match (the default AND behavior).
func NewSearchOr[M any](key Key[M], terms string) Search[M] {
	return &search[M]{key: key, terms: terms, boolOr: true}
}

// Ref sets an explicit predicate reference for this search condition.
// If not called, refs are auto-assigned starting from 0.
func (s *search[M]) Ref(ref int) Search[M] {
	s.ref = &ref
	return s
}

// WithHighlights enables highlighting for this search condition with the given
// prefix and suffix tags.
func (s *search[M]) WithHighlights(prefix, suffix string) Search[M] {
	s.hlPrefix = prefix
	s.hlSuffix = suffix
	s.highlights = true
	return s
}

// WithOffsets enables offset extraction for this search condition.
// Offsets provide the start and end positions of matched terms.
func (s *search[M]) WithOffsets() Search[M] {
	s.offsets = true
	return s
}

func (s *search[M]) build(ctx *context, autoRef int) SearchClause {
	ref := autoRef
	if s.ref != nil {
		ref = *s.ref
	}

	refStr := strconv.Itoa(ref)
	if s.boolOr {
		refStr += ",OR"
	}

	sql := strings.TrimPrefix(s.key.render(ctx), ".") +
		" @" + refStr + "@ '" +
		escapeSearchTerms(s.terms) + "'"

	return SearchClause{
		SQL:        sql,
		Ref:        ref,
		Highlights: s.highlights,
		HLPrefix:   s.hlPrefix,
		HLSuffix:   s.hlSuffix,
		Offsets:    s.offsets,
	}
}

func escapeSearchTerms(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return s
}

// BuildSearchOr combines multiple search conditions with OR and returns the SQL and clauses.
// This is the default behavior used by the Search() method.
// For OR semantics: any condition matching is sufficient.
func BuildSearchOr[M any](searches []Search[M], q *Query[M]) (string, []SearchClause) {
	if len(searches) == 0 {
		return "", nil
	}

	var parts []string
	var clauses []SearchClause
	autoRef := 0

	for _, s := range searches {
		clause := s.build(&q.context, autoRef)
		parts = append(parts, clause.SQL)
		clauses = append(clauses, clause)
		autoRef = clause.Ref + 1
	}

	return "(" + strings.Join(parts, " OR ") + ")", clauses
}

// SearchAll combines multiple search conditions with AND.
// Use this when you need documents to match ALL search terms.
type SearchAll[M any] []Search[M]

// BuildClauses renders all search conditions and returns the SQL and clauses.
func (sa SearchAll[M]) BuildClauses(q *Query[M]) (string, []SearchClause) {
	if len(sa) == 0 {
		return "", nil
	}

	var parts []string
	var clauses []SearchClause
	autoRef := 0

	for _, s := range sa {
		clause := s.build(&q.context, autoRef)
		parts = append(parts, clause.SQL)
		clauses = append(clauses, clause)
		autoRef = clause.Ref + 1
	}

	return "(" + strings.Join(parts, " AND ") + ")", clauses
}
