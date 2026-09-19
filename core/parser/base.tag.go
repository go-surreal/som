package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var validDBName = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// ReservedDBNames contains field names managed by SOM or SurrealDB internals.
var ReservedDBNames = map[string]bool{
	"id":         true,
	"in":         true,
	"out":        true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"expires_at": true,
}

// validExpiryDuration matches SurrealDB duration literals as used in the som.Expiry struct tag,
// e.g. "24h", "7d", "1w", "500ms". Units follow SurrealDB semantics, which are
// broader than Go's time.ParseDuration (it additionally supports d, w and y).
var validExpiryDuration = regexp.MustCompile(`^([0-9]+(ns|us|µs|ms|s|m|h|d|w|y))+$`)

// IndexInfo holds index configuration parsed from struct tags.
type IndexInfo struct {
	// Name is an optional index name from `index=<name>` or `unique=<name>`.
	// For regular indexes, this becomes the SurrealDB index name.
	// For unique indexes with a name, fields sharing the same name are
	// grouped into a single composite unique index.
	// If empty, the index name is auto-generated from table and field names.
	Name string

	// Unique indicates this is a unique index.
	Unique bool
}

// SearchInfo holds fulltext search configuration parsed from struct tags.
type SearchInfo struct {
	// ConfigName references a search configuration defined in a //go:build som file.
	ConfigName string
}

// AssertKind selects how a field constraint is rendered into the schema.
type AssertKind int

const (
	// AssertLen constrains the length of a string or the item count of a slice.
	AssertLen AssertKind = iota
	// AssertNum constrains the range of a numeric value.
	AssertNum
	// AssertConfig references a constraint declared via define.Assert.
	AssertConfig
)

// AssertInfo holds a single field-level constraint parsed from a som struct tag.
type AssertInfo struct {
	Kind AssertKind

	// Min and Max hold the bounds verbatim as written in the tag, so that the
	// schema reproduces the literal the user gave. Either may be empty for a
	// half-open range. Only set for AssertLen and AssertNum.
	Min, Max string

	// ConfigName references a constraint defined in a //go:build som file.
	// Only set for AssertConfig.
	ConfigName string
}

// TagInfo holds all parsed som struct tag data.
type TagInfo struct {
	DBName     string
	Indexes    []IndexInfo
	Search     *SearchInfo
	Asserts    []AssertInfo
	NoValidate bool
}

// validLen matches the bound syntax of the len tag: an exact count ("3") or a
// range with an optionally open end ("3..64", "3..", "..64").
var validLen = regexp.MustCompile(`^(?:([0-9]+)|([0-9]*)\.\.([0-9]*))$`)

// validNumber matches the numeric literals accepted by the min and max tags.
var validNumber = regexp.MustCompile(`^-?[0-9]+(?:\.[0-9]+)?$`)

// parseLenTag turns the value of a len tag into its lower and upper bound.
func parseLenTag(part, value string) (min, max string, err error) {
	match := validLen.FindStringSubmatch(value)
	if match == nil {
		return "", "", fmt.Errorf("invalid tag %q: len expects a count or a range (len=3, len=3..64, len=3.., len=..64)", part)
	}

	if exact := match[1]; exact != "" {
		return exact, exact, nil
	}

	min, max = match[2], match[3]
	if min == "" && max == "" {
		return "", "", fmt.Errorf("invalid tag %q: len range needs at least one bound", part)
	}
	return min, max, nil
}

// parseExpiryTag validates the duration declared on a som.Expiry embed. The whole
// som tag value is the duration (e.g. `som:"24h"`), since the embed already
// conveys the expiry intent. It returns the raw duration string (embedded
// verbatim into the generated schema) or an error.
func parseExpiryTag(tag string) (string, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return "", fmt.Errorf("som.Expiry embed requires a duration via `som:\"<duration>\"` (e.g. som:\"24h\")")
	}
	if !validExpiryDuration.MatchString(tag) {
		return "", fmt.Errorf("som.Expiry embed: %q is not a valid duration", tag)
	}
	return tag, nil
}

// ParseChangefeedTag extracts the changefeed duration from a som tag.
// Tag format: som:"changefeed=2d" returns "2d"
func ParseChangefeedTag(tag string) string {
	if tag == "" {
		return ""
	}

	for part := range strings.SplitSeq(tag, ",") {
		part = strings.TrimSpace(part)
		if duration, ok := strings.CutPrefix(part, "changefeed="); ok {
			return duration
		}
	}
	return ""
}

// setNumBound records a min or max bound on the tag's numeric constraint,
// merging both ends of a range into a single entry.
func setNumBound(info *TagInfo, key, value string) error {
	for i := range info.Asserts {
		assert := &info.Asserts[i]
		if assert.Kind != AssertNum {
			continue
		}
		if key == "min" {
			if assert.Min != "" {
				return fmt.Errorf("min specified multiple times")
			}
			assert.Min = value
		} else {
			if assert.Max != "" {
				return fmt.Errorf("max specified multiple times")
			}
			assert.Max = value
		}
		return nil
	}

	assert := AssertInfo{Kind: AssertNum}
	if key == "min" {
		assert.Min = value
	} else {
		assert.Max = value
	}
	info.Asserts = append(info.Asserts, assert)

	return nil
}

// parseSomTag parses the "som" struct tag and extracts field metadata.
// All parameterized options use key=value syntax:
//
//	som:"index"
//	som:"index=my_index"
//	som:"unique"
//	som:"unique=composite_name"
//	som:"name=db_field_name"
//	som:"fulltext=english_search"
//	som:"index,unique=login"
//	som:"len=3"
//	som:"len=3..64"
//	som:"min=0,max=130"
//	som:"assert=phone_format"
//
// The only flag without a value is:
//
//	som:"novalidate"
func parseSomTag(tag string) (*TagInfo, error) {
	if tag == "" || tag == "in" || tag == "out" {
		return nil, nil
	}

	info := &TagInfo{}

	for part := range strings.SplitSeq(tag, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		key, value, hasValue := strings.Cut(part, "=")

		switch key {
		case "index":
			idx := IndexInfo{}
			if hasValue {
				if value == "" {
					return nil, fmt.Errorf("invalid tag %q: index name must not be empty", part)
				}
				idx.Name = value
			}
			info.Indexes = append(info.Indexes, idx)

		case "unique":
			idx := IndexInfo{Unique: true}
			if hasValue {
				if value == "" {
					return nil, fmt.Errorf("invalid tag %q: unique name must not be empty", part)
				}
				idx.Name = value
			}
			info.Indexes = append(info.Indexes, idx)

		case "name":
			if !hasValue || value == "" {
				return nil, fmt.Errorf("invalid tag %q: name requires a value (name=db_field_name)", part)
			}
			if info.DBName != "" {
				return nil, fmt.Errorf("invalid tag: name specified multiple times")
			}
			if !validDBName.MatchString(value) {
				return nil, fmt.Errorf("invalid tag %q: name must match [a-z][a-z0-9_]*", part)
			}
			if ReservedDBNames[value] {
				return nil, fmt.Errorf("invalid tag %q: %q is a reserved field name", part, value)
			}
			info.DBName = value

		case "fulltext":
			if !hasValue || value == "" {
				return nil, fmt.Errorf("invalid tag %q: fulltext requires a config name (fulltext=english_search)", part)
			}
			info.Search = &SearchInfo{ConfigName: value}

		case "len":
			min, max, err := parseLenTag(part, value)
			if err != nil {
				return nil, err
			}
			info.Asserts = append(info.Asserts, AssertInfo{Kind: AssertLen, Min: min, Max: max})

		case "min", "max":
			if !validNumber.MatchString(value) {
				return nil, fmt.Errorf("invalid tag %q: %s requires a number (%s=10)", part, key, key)
			}
			if err := setNumBound(info, key, value); err != nil {
				return nil, fmt.Errorf("invalid tag %q: %w", part, err)
			}

		case "assert":
			if !hasValue || value == "" {
				return nil, fmt.Errorf("invalid tag %q: assert requires a config name (assert=phone_format)", part)
			}
			info.Asserts = append(info.Asserts, AssertInfo{Kind: AssertConfig, ConfigName: value})

		case "novalidate":
			if hasValue {
				return nil, fmt.Errorf("invalid tag %q: novalidate takes no value", part)
			}
			info.NoValidate = true

		case "changefeed":
			// Handled separately at the node/edge level via ParseChangefeedTag.

		default:
			return nil, fmt.Errorf("unknown som tag %q", part)
		}
	}

	if err := validateAsserts(info.Asserts); err != nil {
		return nil, err
	}

	return info, nil
}

// validateAsserts rejects constraints that can never hold, so that the mistake
// surfaces at generation time instead of on every write.
func validateAsserts(asserts []AssertInfo) error {
	seenLen := false

	for _, assert := range asserts {
		if assert.Kind == AssertLen {
			if seenLen {
				return fmt.Errorf("invalid tag: len specified multiple times")
			}
			seenLen = true
		}

		if assert.Kind == AssertConfig || assert.Min == "" || assert.Max == "" {
			continue
		}

		min, err := strconv.ParseFloat(assert.Min, 64)
		if err != nil {
			return fmt.Errorf("invalid lower bound %q: %w", assert.Min, err)
		}
		max, err := strconv.ParseFloat(assert.Max, 64)
		if err != nil {
			return fmt.Errorf("invalid upper bound %q: %w", assert.Max, err)
		}
		if min > max {
			return fmt.Errorf("invalid tag: lower bound %s is greater than upper bound %s", assert.Min, assert.Max)
		}
	}

	return nil
}
