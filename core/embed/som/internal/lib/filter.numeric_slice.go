//go:build embed

package lib

type NumericSlice[M, E any, F field[M]] struct {
	*Slice[M, E, F]
}

func NewNumericSlice[M, E any, F field[M]](key Key[M], makeElemFilter makeFilter[M, F]) *NumericSlice[M, E, F] {
	return &NumericSlice[M, E, F]{
		Slice: NewSlice[M, E, F](key, makeElemFilter),
	}
}

type NumericSlicePtr[M, E any, F field[M]] struct {
	*SlicePtr[M, E, F]
}

func NewNumericSlicePtr[M, E any, F field[M]](key Key[M], makeElemFilter makeFilter[M, F]) *NumericSlicePtr[M, E, F] {
	return &NumericSlicePtr[M, E, F]{
		SlicePtr: NewSlicePtr[M, E, F](key, makeElemFilter),
	}
}

//
//
//

func (s *NumericSlice[M, E, F]) Interquartile() *Float[M, float64] { // TODO: float or int?
	return NewFloat[M, float64](s.fn("math::interquartile"))
}

func (s *NumericSlice[M, E, F]) Mean() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::mean"))
}

func (s *NumericSlice[M, E, F]) Median() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::median"))
}

func (s *NumericSlice[M, E, F]) Midhinge() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::midhinge"))
}

func (s *NumericSlice[M, E, F]) NearestRank() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::nearestrank"))
}

func (s *NumericSlice[M, E, F]) Percentile() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::percentile"))
}

func (s *NumericSlice[M, E, F]) StdDev() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::stddev"))
}

func (s *NumericSlice[M, E, F]) TriMean() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::trimean"))
}

func (s *NumericSlice[M, E, F]) Variance() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::variance"))
}
