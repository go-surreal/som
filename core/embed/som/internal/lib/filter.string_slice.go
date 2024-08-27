//go:build embed

package lib

type StringSlice[M any] struct {
	*Slice[M, string, *String[M]]
}

func NewStringSlice[M any](key Key[M]) *StringSlice[M] {
	return &StringSlice[M]{
		Slice: NewSlice[M, string, *String[M]](key, NewString[M]),
	}
}

type StringSlicePtr[M any] struct {
	*SlicePtr[M, string, *String[M]]
}

func NewStringSlicePtr[M any](key Key[M]) *StringSlicePtr[M] {
	return &StringSlicePtr[M]{
		SlicePtr: NewSlicePtr[M, string, *String[M]](key, NewString[M]),
	}
}

type StringPtrSlice[M any] struct {
	*Slice[M, *string, *StringPtr[M]]
}

func NewStringPtrSlice[M any](key Key[M]) *StringPtrSlice[M] {
	return &StringPtrSlice[M]{
		Slice: NewSlice[M, *string, *StringPtr[M]](key, NewStringPtr[M]),
	}
}

type StringPtrSlicePtr[M any] struct {
	*SlicePtr[M, *string, *StringPtr[M]]
}

func NewStringPtrSlicePtr[M any](key Key[M]) *StringPtrSlicePtr[M] {
	return &StringPtrSlicePtr[M]{
		SlicePtr: NewSlicePtr[M, *string, *StringPtr[M]](key, NewStringPtr[M]),
	}
}

//
//
//

func (s *StringSlice[M]) Join(sep string) *String[M] {
	return NewString(s.fn("string::join", sep))
}
