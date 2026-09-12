package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// Invalid matches a field whose type is a model type that exists in the
// package but can never be referenced from another model. Only nodes and edges
// may be referenced, so such a field is always invalid.
//
// Such a field would otherwise fall through to FieldStruct, whose Verify
// rejects it as an unknown struct — the referenced type is not unknown at all,
// and Verify only runs after the (slow) definition step. Matching the invalid
// reference here fails immediately and names the actual reason.
//
// It produces no field and therefore implements FieldHandler only.
type Invalid struct {
	kind   string
	detect func(t gotype.Type, outPkg string) bool
	reason string
}

// InvalidView marks a reference to a view as invalid. Views are read-only,
// computed tables; their rows cannot be linked to, so a record link (or any
// field) pointing at a view is not allowed.
func InvalidView() *Invalid {
	return &Invalid{
		kind:   "view",
		detect: IsView,
		reason: "views are read-only and cannot be linked",
	}
}

// InvalidSink marks a reference to a sink as invalid. Sink records are
// discarded immediately after write, so they have no addressable id.
func InvalidSink() *Invalid {
	return &Invalid{
		kind:   "sink",
		detect: IsSink,
		reason: "sink records are discarded after write and cannot be linked",
	}
}

// InvalidFragment marks a reference to a fragment as invalid. A fragment is
// a projection of a node, not a table of its own, so it cannot be stored.
func InvalidFragment() *Invalid {
	return &Invalid{
		kind:   "fragment",
		detect: IsFragment,
		reason: "a fragment is a projection of a node and cannot be stored or linked",
	}
}

func (h *Invalid) Match(elem gotype.Type, ctx *FieldContext) bool {
	return h.detect(elem, ctx.OutPkg)
}

func (h *Invalid) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return nil, fmt.Errorf(
		"field %q references %s %q: %s (only nodes and edges may be referenced)",
		t.Name(), h.kind, elem.Name(), h.reason,
	)
}
