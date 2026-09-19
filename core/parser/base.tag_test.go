package parser

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseSomTagAsserts(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		asserts []AssertInfo
	}{
		{
			name:    "exact length",
			tag:     "len=4",
			asserts: []AssertInfo{{Kind: AssertLen, Min: "4", Max: "4"}},
		},
		{
			name:    "length range",
			tag:     "len=3..64",
			asserts: []AssertInfo{{Kind: AssertLen, Min: "3", Max: "64"}},
		},
		{
			name:    "length lower bound only",
			tag:     "len=3..",
			asserts: []AssertInfo{{Kind: AssertLen, Min: "3"}},
		},
		{
			name:    "length upper bound only",
			tag:     "len=..64",
			asserts: []AssertInfo{{Kind: AssertLen, Max: "64"}},
		},
		{
			name:    "numeric bounds merge into one constraint",
			tag:     "min=0,max=130",
			asserts: []AssertInfo{{Kind: AssertNum, Min: "0", Max: "130"}},
		},
		{
			name:    "negative and fractional bounds",
			tag:     "min=-1.5,max=1.5",
			asserts: []AssertInfo{{Kind: AssertNum, Min: "-1.5", Max: "1.5"}},
		},
		{
			name:    "lower bound only",
			tag:     "min=0",
			asserts: []AssertInfo{{Kind: AssertNum, Min: "0"}},
		},
		{
			name:    "named config",
			tag:     "assert=phone_format",
			asserts: []AssertInfo{{Kind: AssertConfig, ConfigName: "phone_format"}},
		},
		{
			name: "combined with other options",
			tag:  "index,len=1..80,assert=no_placeholder",
			asserts: []AssertInfo{
				{Kind: AssertLen, Min: "1", Max: "80"},
				{Kind: AssertConfig, ConfigName: "no_placeholder"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info, err := parseSomTag(test.tag)
			assert.NilError(t, err)
			assert.DeepEqual(t, test.asserts, info.Asserts)
		})
	}
}

func TestParseSomTagAssertErrors(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{name: "empty length", tag: "len="},
		{name: "open range", tag: "len=.."},
		{name: "non numeric length", tag: "len=abc"},
		{name: "negative length", tag: "len=-1"},
		{name: "reversed length range", tag: "len=64..3"},
		{name: "reversed numeric range", tag: "min=10,max=1"},
		{name: "non numeric bound", tag: "min=abc"},
		{name: "duplicate length", tag: "len=3,len=4"},
		{name: "duplicate minimum", tag: "min=1,min=2"},
		{name: "duplicate maximum", tag: "max=1,max=2"},
		{name: "empty config name", tag: "assert="},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseSomTag(test.tag)
			assert.Assert(t, err != nil, "expected %q to be rejected", test.tag)
		})
	}
}
