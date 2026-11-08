//go:build embed

package lib

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	somPrefix             = "__som__"
	searchScorePrefix     = somPrefix + "search_score_"
	searchHighlightPrefix = somPrefix + "search_highlight_"
	searchOffsetsPrefix   = somPrefix + "search_offsets_"
)

type context struct {
	varIndex int32
	vars     map[string]any

	// literal makes asVar inline values as SurrealQL literals instead of
	// binding them as $-parameters. This is required when rendering static
	// DDL such as a DEFINE TABLE ... AS SELECT view definition, which cannot
	// carry query parameters.
	literal bool
}

func (c *context) Vars() map[string]any {
	return c.vars
}

// TODO: deduplicate same values in the same query
func (c *context) asVar(val any) string {
	if c.literal {
		return literalValue(val)
	}

	index := intToLetters(c.varIndex)
	c.varIndex++
	c.vars[index] = val
	return "$" + index
}

type Query[T any] struct {
	context
	node       string
	live       bool
	fields     []string
	groupBy    string
	groupAll   bool
	Where      []Filter[T]
	Sort       []*SortBuilder
	SortRandom bool
	Fetch      []string
	Start      int
	Limit      int
	Timeout    time.Duration
	TempFiles  bool

	// Range restricts the query to a range of record IDs.
	Range *RecordRange

	// SeekAfter holds the query variable of a cursor record ID. When set, the
	// query scans the record ID range exclusively after it instead of filtering
	// the whole table, which is only valid for an "id ASC" order.
	SeekAfter string

	SearchClauses []SearchClause
	SearchWhere   string

	// Soft delete support for main queries (not fetched relations)
	SoftDeleteFilter Filter[T] // Injected at initialization
	IncludeDeleted   bool      // Flag to skip soft delete filter

	// Expiry (TTL) support: when set, expired records are excluded from main queries.
	ExpiryField    string // Database field name holding the expiry timestamp
	IncludeExpired bool   // Flag to skip the expiry filter
}

func (q *Query[T]) AsVar(val any) string {
	return q.context.asVar(val)
}

func NewQuery[T any](node string) Query[T] {
	return Query[T]{
		context: context{
			varIndex: 0,
			vars:     map[string]any{},
		},
		node: node,
	}
}

func (q Query[T]) BuildAsAll() *Result {
	q.fields = []string{"*"}

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildAsAllIDs() *Result {
	q.fields = []string{"id"}

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

// BuildAsFragment builds the query with the projection narrowed to the given
// database fields, as used for fragment queries.
func (q Query[T]) BuildAsFragment(fields []string) *Result {
	q.fields = q.projection(fields)

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

// BuildAsLiveFragment is the live version of BuildAsFragment.
func (q Query[T]) BuildAsLiveFragment(fields []string) *Result {
	q.live = true
	q.fields = q.projection(fields)

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

// projection returns the select list for a narrowed query: the requested
// fields plus every sort field that is missing from them. ORDER BY resolves
// against the projected record, so a sort on an unselected field would have
// nothing to order by.
func (q Query[T]) projection(fields []string) []string {
	out := slices.Clone(fields)

	for _, s := range q.Sort {
		if s == nil || s.Field == "" || s.IsScore {
			continue
		}
		if !slices.Contains(out, s.Field) {
			out = append(out, s.Field)
		}
	}

	return out
}

func (q Query[T]) BuildAsCount() *Result {
	q.fields = []string{"count()"}
	q.groupAll = true

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildAsLive() *Result {
	q.live = true
	q.fields = []string{"*"}

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildAsLiveDiff() *Result {
	q.live = true
	q.fields = []string{"DIFF"}

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildDistinct(field string) *Result {
	q.fields = []string{"array::distinct(array::group(" + field + ")) AS values"}
	q.groupAll = true

	q.Where = append(q.Where, filter[T](func(_ *context, _ T) string {
		return "(" + field + " != NONE AND " + field + " != NULL)"
	}))

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) render() string {
	var out strings.Builder

	// TODO: possible optimization: preallocate buffer (e.g. out.Grow(<known bytes>))

	if q.live {
		out.WriteString("LIVE ")
	}

	fields := q.fields

	// Add search score, highlight, and offset projections
	for _, sc := range q.SearchClauses {
		ref := strconv.Itoa(sc.Ref)
		fields = append(fields, "search::score("+ref+") AS "+searchScorePrefix+ref)
		if sc.Highlights {
			fields = append(fields,
				"search::highlight("+q.context.asVar(sc.HLPrefix)+", "+q.context.asVar(sc.HLSuffix)+", "+ref+") AS "+searchHighlightPrefix+ref)
		}
		if sc.Offsets {
			fields = append(fields, "search::offsets("+ref+") AS "+searchOffsetsPrefix+ref)
		}
	}

	// Add score projection for each score sort
	for _, s := range q.Sort {
		if s.IsScore && len(s.ScoreRefs) > 0 {
			expr := renderScoreCombination(s.ScoreRefs, s.ScoreMode, s.ScoreWeights)
			refStrs := make([]string, len(s.ScoreRefs))
			for i, ref := range s.ScoreRefs {
				refStrs[i] = strconv.Itoa(ref)
			}
			alias := searchScorePrefix + strings.Join(refStrs, "_")
			fields = append(fields, expr+" AS "+alias)
		}
	}

	out.WriteString("SELECT " + strings.Join(fields, ", "))

	// TODO: not working, but more a optimisation than anything else
	//if len(q.Sort) > 0 {
	//	out.WriteString(" OMIT ")
	//	var omitFields []string
	//	for _, s := range q.Sort {
	//		omitFields = append(omitFields, "__som_"+s.Field)
	//	}
	//	out.WriteString(strings.Join(omitFields, ", "))
	//}

	out.WriteString(" FROM " + q.source())

	// Build WHERE clause combining soft delete, search and regular filters
	var whereParts []string

	// 1. Inject soft delete filter FIRST (if enabled and not disabled)
	if !q.IncludeDeleted && q.SoftDeleteFilter != nil {
		var t T
		if sdFilter := q.SoftDeleteFilter.build(&q.context, t); sdFilter != "" {
			whereParts = append(whereParts, sdFilter)
		}
	}

	// 1b. Inject expiry filter to exclude expired records (if enabled and not disabled)
	if !q.IncludeExpired && q.ExpiryField != "" {
		whereParts = append(whereParts, "("+q.ExpiryField+" IS NONE OR "+q.ExpiryField+" > time::now())")
	}

	// 1c. Keep the lower bound of a range that a cursor seek replaced
	if q.SeekAfter != "" && q.Range != nil {
		if guard := q.Range.renderFromGuard(q.node); guard != "" {
			whereParts = append(whereParts, guard)
		}
	}

	// 2. Add search conditions
	if q.SearchWhere != "" {
		whereParts = append(whereParts, q.SearchWhere)
	}

	// 3. Add regular filters
	var t T
	filterWhere := All[T](q.Where).build(&q.context, t)
	if filterWhere != "" {
		whereParts = append(whereParts, filterWhere)
	}

	if len(whereParts) > 0 {
		out.WriteString(" WHERE ")
		out.WriteString(strings.Join(whereParts, " AND "))
	}

	if !q.live && q.groupBy != "" {
		out.WriteString(" GROUP BY ")
		out.WriteString(q.groupBy)
	}

	if !q.live && q.groupAll {
		out.WriteString(" GROUP ALL")
	}

	if !q.live && q.SortRandom {
		out.WriteString(" ORDER BY RAND()")
	} else if !q.live && len(q.Sort) > 0 {
		out.WriteString(" ORDER BY ")
		out.WriteString(q.renderSorts())
	}

	// LIMIT must come before START.
	if !q.live && q.Limit > 0 {
		out.WriteString(" LIMIT ")
		out.WriteString(strconv.Itoa(q.Limit))
	}

	// START must come after LIMIT.
	if !q.live && q.Start > 0 {
		out.WriteString(" START ")
		out.WriteString(strconv.Itoa(q.Start))
	}

	if len(q.Fetch) > 0 {
		out.WriteString(" FETCH ")
		// Note: Soft-delete filtering does NOT apply to fetched relations.
		// All related records are returned regardless of their soft-delete status.
		// Users should filter manually using IsDeleted() if needed.
		out.WriteString(strings.Join(q.Fetch, ", "))
	}

	if !q.live && q.Timeout > 0 {
		out.WriteString(" TIMEOUT ")
		out.WriteString(q.Timeout.Round(time.Second).String())
	}

	if !q.live && q.TempFiles {
		out.WriteString(" TEMPFILES")
	}

	return out.String()
}

// source returns the FROM target of the query. Record ID ranges are wrapped in
// type::record() because SurrealDB rejects params inside a literal range key
// (table:$a..$b), and because only this dynamic form pushes LIMIT into the
// range scan instead of reading the whole range first.
func (q Query[T]) source() string {
	switch {
	case q.SeekAfter != "":
		return recordFromKey(q.node, q.recordRange().seekAfter(q.SeekAfter).render())

	case q.Range != nil:
		return recordFromKey(q.node, q.Range.render())

	default:
		return q.node
	}
}

func (q Query[T]) recordRange() RecordRange {
	if q.Range == nil {
		return RecordRange{}
	}

	return *q.Range
}

func (q Query[T]) renderSorts() string {
	sorts := make([]string, 0, len(q.Sort)+1)

	for _, sort := range q.Sort {
		sorts = append(sorts, sort.render())
	}

	if q.isSoleRangeIDDescSort() {
		sorts = append(sorts, sorts[0])
	}

	return strings.Join(sorts, ", ")
}

// isSoleRangeIDDescSort reports whether the query orders a record ID range by
// "id DESC" alone. SurrealDB (3.3.5) pushes LIMIT into such a type::record()
// range scan before sorting, returning the lowest IDs instead of the highest.
// Repeating the sort key disables that pushdown.
func (q Query[T]) isSoleRangeIDDescSort() bool {
	if q.Range == nil || len(q.Sort) != 1 {
		return false
	}

	sort := q.Sort[0]

	return sort.Field == "id" && sort.Order == SortDesc
}

type Result struct {
	Statement string
	Variables map[string]any
}

type Operator string

const (
	OpAnd Operator = "AND" // ("&&") Checks whether both of two values are truthy.
	OpOr  Operator = "OR"  // ("||") Checks whether either of two values is truthy.

	OpEqual         Operator = "="  // ("IS") Check whether two values are equal.
	OpNotEqual      Operator = "!=" // ("IS NOT") Check whether two values are not equal.
	OpExactlyEqual  Operator = "==" // Check whether two values are exactly equal.
	OpFuzzyMatch    Operator = "~"  // Compare two values for equality using fuzzy matching.
	OpFuzzyNotMatch Operator = "!~" // Compare two values for inequality using fuzzy matching.

	OpAnyEqual      Operator = "?=" // Check whether any value in a set is equal to a value.
	OpAllEqual      Operator = "*=" // Check whether all values in a set are equal to a value.
	OpAnyFuzzyMatch Operator = "?~" // Check whether any value in a set is equal to a value using fuzzy matching.
	OpAllFuzzyMatch Operator = "*~" // Check whether all values in a set are equal to a value using fuzzy matching.

	OpLessThan         Operator = "<"  // Check whether a value is less than another value.
	OpLessThanEqual    Operator = "<=" // Check whether a value is less than or equal to another value.
	OpGreaterThan      Operator = ">"  // Check whether a value is greater than another value.
	OpGreaterThanEqual Operator = ">=" // Check whether a value is greater than or equal to another value.

	OpAdd   Operator = "+"  // 	Add two values together.
	OpSub   Operator = "-"  // Subtract a value from another value.
	OpMul   Operator = "×"  // ("*") Multiply two values together.
	OpDiv   Operator = "÷"  // ("/") Divide a value by another value.
	OpRaise Operator = "**" // Raises a base value by another value.

	OpInvert               Operator = "!"  // Reverses the truthiness of a value.
	OpTruth                Operator = "!!" // Determines the truthiness of a value.
	OpEitherTrueAndNotNull Operator = "??" // Check whether either of two values are truthy and not NULL.
	OpEitherTrue           Operator = "?:" // Check whether either of two values are truthy.

	OpContains     Operator = "∋" // ("CONTAINS") Checks whether a value contains another value.
	OpContainsNot  Operator = "∌" // ("CONTAINSNOT") Checks whether a value does not contain another value.
	OpContainsAll  Operator = "⊇" // ("CONTAINSALL") Checks whether a value contains all other values.
	OpContainsAny  Operator = "⊃" // ("CONTAINSANY") Checks whether a value contains any other value.
	OpContainsNone Operator = "⊅" // ("CONTAINSNONE") Checks whether a value contains none of the following values.

	OpIn     Operator = "∈" // ("INSIDE") Checks whether a value is contained within another value. - TODO: geo!
	OpNotIn  Operator = "∉" // ("NOTINSIDE" | "NOT IN") Checks whether a value is not contained within another value. - TODO: geo!
	OpAllIn  Operator = "⊆" // ("ALLINSIDE") Checks whether all values are contained within other values.
	OpAnyIn  Operator = "⊂" // ("ANYINSIDE") Checks whether any value is contained within other values.
	OpNoneIn Operator = "⊄" // ("NONEINSIDE") Checks whether no value is contained within other values.

	OpGeoOutside    Operator = "OUTSIDE"    // Checks whether a geometry type is outside another geometry type.
	OpGeoIntersects Operator = "INTERSECTS" // Checks whether a geometry type intersects another geometry type.

	OpSearch Operator = "@@" // ("@[ref]@") Checks whether the terms are found in a full-text indexed field. - TODO!

	OpX = "<|4|> or <|3,HAMMING|>" // KNN - TODO!

	CastInt   Operator = "<int>"
	CastFloat Operator = "<float>"

	OpModulo Operator = "%" // https://github.com/surrealdb/surrealdb/pull/4182
)

func renderScoreCombination(refs []int, mode ScoreCombineMode, weights []float64) string {
	var scoreParts []string
	for _, ref := range refs {
		scoreParts = append(scoreParts, "search::score("+strconv.Itoa(ref)+")")
	}

	switch mode {
	case ScoreCombineMax:
		return "math::max(" + strings.Join(scoreParts, ", ") + ")"
	case ScoreCombineAverage:
		return "((" + strings.Join(scoreParts, " + ") + ") / " + strconv.Itoa(len(refs)) + ")"
	case ScoreCombineWeighted:
		var weightedParts []string
		for i, ref := range refs {
			weightedParts = append(weightedParts,
				"search::score("+strconv.Itoa(ref)+") * "+strconv.FormatFloat(weights[i], 'f', -1, 64))
		}
		return "(" + strings.Join(weightedParts, " + ") + ")"
	default: // ScoreCombineSum
		return "(" + strings.Join(scoreParts, " + ") + ")"
	}
}

//
// -- HELPER
//

var letterDef = [...]string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

func intToLetters(number int32) (letters string) {
	if firstLetter := number / 26; firstLetter > 0 {
		letters += intToLetters(firstLetter)
		letters += letterDef[number%26]
		return
	}

	letters += letterDef[number]
	return
}
