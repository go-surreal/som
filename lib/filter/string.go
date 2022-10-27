package filter

type String[R any] struct {
	*Base[string, R]
	*Comparable[string, R]
}

func NewString[R any](key string) *String[R] {
	return &String[R]{
		Base:       &Base[string, R]{key: key},
		Comparable: &Comparable[string, R]{key: key},
	}
}

func (String[R]) FuzzyMatch(val string) *Of[R] {
	return &Of[R]{}
}

func (String[R]) NotFuzzyMatch(val string) *Of[R] {
	return &Of[R]{}
}
