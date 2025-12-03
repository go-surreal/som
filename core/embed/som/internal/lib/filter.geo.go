//go:build embed

package lib

// https://surrealdb.com/docs/surrealdb/surrealql/functions/database/geo

// Geo is a base filter for geometry types with spatial operations.
type Geo[M, T any, F field[M]] struct {
	*Base[M, T, F, *Slice[M, T, F]]
}

func NewGeo[M, T any, F field[M]](key Key[M], conv func(T) any) *Geo[M, T, F] {
	return &Geo[M, T, F]{
		Base: NewBaseConv[M, T, F, *Slice[M, T, F]](key, conv),
	}
}

// Contains checks if this geometry contains the other geometry.
func (g *Geo[M, T, F]) Contains(other T) Filter[M] {
	if g.conv != nil {
		return g.Key.op(OpContains, g.conv(other))
	}
	return g.Key.op(OpContains, other)
}

// Inside checks if this geometry is inside the other geometry.
func (g *Geo[M, T, F]) Inside(other T) Filter[M] {
	if g.conv != nil {
		return g.Key.op(OpIn, g.conv(other))
	}
	return g.Key.op(OpIn, other)
}

// Outside checks if this geometry is outside the other geometry.
func (g *Geo[M, T, F]) Outside(other T) Filter[M] {
	if g.conv != nil {
		return g.Key.op(OpGeoOutside, g.conv(other))
	}
	return g.Key.op(OpGeoOutside, other)
}

// Intersects checks if this geometry intersects with the other geometry.
func (g *Geo[M, T, F]) Intersects(other T) Filter[M] {
	if g.conv != nil {
		return g.Key.op(OpGeoIntersects, g.conv(other))
	}
	return g.Key.op(OpGeoIntersects, other)
}

// GeoPtr is a pointer version of Geo that adds nil checks.
type GeoPtr[M, T any, F field[M]] struct {
	*Geo[M, T, F]
	*Nillable[M]
}

func NewGeoPtr[M, T any, F field[M]](key Key[M], conv func(T) any) *GeoPtr[M, T, F] {
	return &GeoPtr[M, T, F]{
		Geo:      NewGeo[M, T, F](key, conv),
		Nillable: NewNillable[M](key),
	}
}
