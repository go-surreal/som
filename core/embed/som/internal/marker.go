//go:build embed

package internal

// Marker is a bit set describing the load state of a model instance.
type Marker uint32

const (
	// MarkerLoaded is set on every model instance that originates from a
	// database record.
	MarkerLoaded Marker = 1 << iota

	// MarkerPartial is set on model instances that do not hold all fields of
	// the underlying record and must therefore not be written back to the
	// database. Currently the only such instances are record links that were
	// not resolved via Fetch, which carry just their record id. Future
	// partial loads (e.g. field projections) will use the same flag.
	MarkerPartial

	// MarkerDeleted is set on model instances that were permanently deleted
	// from the database. Writing them back would recreate the record, so all
	// further write operations on them are rejected.
	MarkerDeleted

	// MarkerFromCache is set on model instances that are held by an in-process
	// cache. Such an instance is shared by every reader of that cache entry,
	// so it must not be mutated and may be stale. It is purely informational.
	MarkerFromCache
)

// Has reports whether all of the given flags are set.
func (m Marker) Has(flags Marker) bool {
	return m&flags == flags
}

// LoadMarker holds the load state of a model instance. It is embedded into
// som.Node, som.Edge and som.View under an unexported name, so that model
// instances expose the state read-only, while only the generated code can
// change it.
type LoadMarker struct {
	flags Marker

	// fetched holds one bit per relation field of the model, set when that
	// field is known to hold no unresolved links. The bit order is defined by
	// the generated code that owns the model.
	fetched uint64
}

// Marker returns the load state of the model instance.
func (m LoadMarker) Marker() Marker {
	return m.flags
}

// Fetched returns the relation fields of the model instance that hold no
// unresolved links, as a bit set.
func (m LoadMarker) Fetched() uint64 {
	return m.fetched
}

// IsPartial reports whether the model instance was only partially loaded,
// e.g. as a record link that was not resolved via Fetch. Partial instances
// hold their id, but no field values, so they cannot be written back.
func (m LoadMarker) IsPartial() bool {
	return m.flags.Has(MarkerPartial)
}

func (m *LoadMarker) setMarker(flags Marker) {
	m.flags = flags
}

func (m *LoadMarker) addMarker(flags Marker) {
	m.flags |= flags
}

func (m *LoadMarker) addFetched(bits uint64) {
	m.fetched |= bits
}

type markable interface {
	setMarker(Marker)
	addMarker(Marker)
	addFetched(uint64)
	Marker() Marker
	Fetched() uint64
}

// SetMarker replaces the load state of a node, edge or view.
func SetMarker(target markable, flags Marker) {
	target.setMarker(flags)
}

// AddMarker adds the given flags to the load state of a node, edge or view.
func AddMarker(target markable, flags Marker) {
	target.addMarker(flags)
}

// AddMarkerAny adds the given flags to the load state of v, if v is a node,
// edge or view. It exists for generic code that cannot name the model type.
func AddMarkerAny(v any, flags Marker) {
	if target, ok := v.(markable); ok {
		target.addMarker(flags)
	}
}

// AddFetched flags the given relation fields of a node, edge or view as
// holding no unresolved links.
func AddFetched(target markable, bits uint64) {
	target.addFetched(bits)
}

// Fetched returns the resolved relation fields of v, if v is a node, edge or
// view. It exists for generic code that cannot name the model type.
func Fetched(v any) uint64 {
	if target, ok := v.(markable); ok {
		return target.Fetched()
	}
	return 0
}

// AddFetchedAny flags the given relation fields of v as holding no unresolved
// links, if v is a node, edge or view.
func AddFetchedAny(v any, bits uint64) {
	if target, ok := v.(markable); ok {
		target.addFetched(bits)
	}
}

// MarkerAny returns the load state of v, if v is a node, edge or view.
func MarkerAny(v any) Marker {
	if target, ok := v.(markable); ok {
		return target.Marker()
	}
	return 0
}

// SetMarkerAny replaces the load state of v, if v is a node, edge or view.
func SetMarkerAny(v any, flags Marker) {
	if target, ok := v.(markable); ok {
		target.setMarker(flags)
	}
}
