//go:build embed

package lib

// String is a filter builder for string values.
// M is the model this filter is for.
type String[M any] struct {
	*Base[M, string]
	*Comparable[M, string]
}

func NewString[M any](key Key[M]) *String[M] {
	return &String[M]{
		Base:       NewBase[M, string](key),
		Comparable: NewComparable[M, string](key),
	}
}

func (s *String[M]) Equal_(val *String[M]) Filter[M] {
	return s.Base.key.op_(OpEqual, val.Base.key)
}

func (s *String[M]) FuzzyMatch(val string) Filter[M] {
	return s.Base.key.op(OpFuzzyMatch, val)
}

func (s *String[M]) NotFuzzyMatch(val string) Filter[M] {
	return s.Base.key.op(OpFuzzyNotMatch, val)
}

func (s *String[M]) Contains(val string) *Bool[M] {
	return NewBool(s.Base.key.fn("string::contains", val))
}

func (s *String[M]) EndsWith(val string) *Bool[M] {
	return NewBool(s.Base.key.fn("string::endsWith", val))
}

func (s *String[M]) Len() *Numeric[M, int] {
	return NewNumeric[M, int](s.Base.key.fn("string::len"))
}

func (s *String[M]) Lowercase() *String[M] {
	return NewString(s.Base.key.fn("string::lowercase"))
}

func (s *String[M]) Reverse() *String[M] {
	return NewString(s.Base.key.fn("string::reverse"))
}

func (s *String[M]) Slice(start, end int) *String[M] {
	return NewString(s.Base.key.fn("string::slice", start, end))
}

func (s *String[M]) Slug() *String[M] {
	return NewString(s.Base.key.fn("string::slug"))
}

func (s *String[M]) Split(sep string) *Slice[M, string, *String[M]] {
	return NewSlice[M, string, *String[M]](s.Base.key.fn("string::split", sep), NewString[M])
}

func (s *String[M]) StartsWith(val string) *Bool[M] {
	return NewBool(s.Base.key.fn("string::startsWith", val))
}

func (s *String[M]) Trim() *String[M] {
	return NewString(s.Base.key.fn("string::trim"))
}

func (s *String[M]) Uppercase() *String[M] {
	return NewString(s.Base.key.fn("string::uppercase"))
}

func (s *String[M]) Words() *Slice[M, string, *String[M]] {
	return NewSlice[M, string, *String[M]](s.Base.key.fn("string::words"), NewString[M])
}

func (s *String[M]) IsAlphaNum() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::alphanum"))
}

func (s *String[M]) IsAlpha() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::alpha"))
}

func (s *String[M]) IsAscii() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::ascii"))
}

func (s *String[M]) IsDateTime() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::datetime"))
}

func (s *String[M]) IsDomain() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::domain"))
}

func (s *String[M]) IsEmail() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::email"))
}

func (s *String[M]) IsHexadecimal() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::hexadecimal"))
}

func (s *String[M]) IsIP() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::ip"))
}

func (s *String[M]) IsIPv4() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::ipv4"))
}

func (s *String[M]) IsIPv6() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::ipv6"))
}

func (s *String[M]) IsLatitude() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::latitude"))
}

func (s *String[M]) IsLongitude() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::longitude"))
}

func (s *String[M]) IsNumeric() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::numeric"))
}

func (s *String[M]) IsSemVer() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::semver"))
}

func (s *String[M]) IsURL() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::url"))
}

func (s *String[M]) IsUUID() *Bool[M] {
	return NewBool(s.Base.key.fn("string::is::uuid"))
}

func (s *String[M]) Base64Decode() *ByteSlice[M] {
	return NewByteSlice(s.Base.key.fn("encoding::base64::decode"))
}

type StringPtr[M any] struct {
	*String[M]
	*Nillable[M]
}

func NewStringPtr[M any](key Key[M]) *StringPtr[M] {
	return &StringPtr[M]{
		String:   NewString[M](key),
		Nillable: NewNillable[M](key),
	}
}
