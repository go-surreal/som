package parser

import (
	"fmt"
	"regexp"
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

// TagInfo holds all parsed som struct tag data.
type TagInfo struct {
	DBName  string
	Indexes []IndexInfo
	Search  *SearchInfo
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

		case "changefeed":
			// Handled separately at the node/edge level via ParseChangefeedTag.

		default:
			return nil, fmt.Errorf("unknown som tag %q", part)
		}
	}

	return info, nil
}
