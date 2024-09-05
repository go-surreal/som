//go:build embed

package lib

// String is a filter builder for string values.
// M is the model this filter is for.
type String[M any] struct {
	*Base[M, string, *String[M], *Slice[M, string, *String[M]]]
	*Comparable[M, string, *String[M]]
}

func NewString[M any](key Key[M]) *String[M] {
	return &String[M]{
		Base:       NewBase[M, string, *String[M], *Slice[M, string, *String[M]]](key),
		Comparable: NewComparable[M, string, *String[M]](key),
	}
}

func (s *String[M]) key() Key[M] {
	return s.Base.key()
}

func (s *String[M]) FuzzyMatch(val string) Filter[M] {
	return s.Base.op(OpFuzzyMatch, val)
}

func (s *String[M]) FuzzyNotMatch(val string) Filter[M] {
	return s.Base.op(OpFuzzyNotMatch, val)
}

func (s *String[M]) Concat(vals ...string) *Bool[M] {
	anys := make([]any, len(vals))

	for i, val := range vals {
		anys[i] = val
	}

	return NewBool(s.Base.fn("string::concat", anys...))
}

func (s *String[M]) Contains(val string) *Bool[M] {
	return NewBool(s.Base.fn("string::contains", val)) // or: WHERE val IN key
}

func (s *String[M]) EndsWith(val string) *Bool[M] {
	return NewBool(s.Base.fn("string::endsWith", val))
}

// Join joins the given strings with the base string as separator.
func (s *String[M]) Join(vals ...string) *String[M] {
	anys := make([]any, len(vals))

	for i, val := range vals {
		anys[i] = val
	}

	return NewString(s.Base.fn("string::join", anys...))
}

func (s *String[M]) Len() *Numeric[M, int] {
	return NewNumeric[M, int](s.Base.fn("string::len"))
}

func (s *String[M]) Lowercase() *String[M] {
	return NewString(s.Base.fn("string::lowercase"))
}

func (s *String[M]) Repeat(times int) *String[M] {
	return NewString(s.Base.fn("string::repeat", times))
}

func (s *String[M]) Replace(old, new string) *String[M] {
	return NewString(s.Base.fn("string::replace", old, new))
}

func (s *String[M]) Reverse() *String[M] {
	return NewString(s.Base.fn("string::reverse"))
}

func (s *String[M]) Slice(start, end int) *String[M] {
	return NewString(s.Base.fn("string::slice", start, end))
}

func (s *String[M]) Slug() *String[M] {
	return NewString(s.Base.fn("string::slug"))
}

func (s *String[M]) Split(sep string) *Slice[M, string, *String[M]] {
	return NewSlice[M, string, *String[M]](s.Base.fn("string::split", sep), NewString[M])
}

func (s *String[M]) StartsWith(val string) *Bool[M] {
	return NewBool(s.Base.fn("string::startsWith", val))
}

func (s *String[M]) Trim() *String[M] {
	return NewString(s.Base.fn("string::trim"))
}

func (s *String[M]) Uppercase() *String[M] {
	return NewString(s.Base.fn("string::uppercase"))
}

func (s *String[M]) Words() *Slice[M, string, *String[M]] {
	return NewSlice[M, string, *String[M]](s.Base.fn("string::words"), NewString[M])
}

// TODO: string::html::encode (v2.0.0)
// TODO: string::html::sanitize (v2.0.0)

func (s *String[M]) IsAlphaNum() *Bool[M] {
	return NewBool(s.Base.fn("string::is::alphanum"))
}

func (s *String[M]) IsAlpha() *Bool[M] {
	return NewBool(s.Base.fn("string::is::alpha"))
}

func (s *String[M]) IsAscii() *Bool[M] {
	return NewBool(s.Base.fn("string::is::ascii"))
}

// IsDateTime returns a boolean field that is true if the string is a valid date and time.
//
// See [Time.Format] for the format string.
func (s *String[M]) IsDateTime(format string) *Bool[M] {
	return NewBool(s.Base.fn("string::is::datetime", format))
}

func (s *String[M]) IsDomain() *Bool[M] {
	return NewBool(s.Base.fn("string::is::domain"))
}

func (s *String[M]) IsEmail() *Bool[M] {
	return NewBool(s.Base.fn("string::is::email"))
}

func (s *String[M]) IsHexadecimal() *Bool[M] {
	return NewBool(s.Base.fn("string::is::hexadecimal"))
}

func (s *String[M]) IsIP() *Bool[M] {
	return NewBool(s.Base.fn("string::is::ip"))
}

func (s *String[M]) IsIPv4() *Bool[M] {
	return NewBool(s.Base.fn("string::is::ipv4"))
}

func (s *String[M]) IsIPv6() *Bool[M] {
	return NewBool(s.Base.fn("string::is::ipv6"))
}

func (s *String[M]) IsLatitude() *Bool[M] {
	return NewBool(s.Base.fn("string::is::latitude"))
}

func (s *String[M]) IsLongitude() *Bool[M] {
	return NewBool(s.Base.fn("string::is::longitude"))
}

func (s *String[M]) IsNumeric() *Bool[M] {
	return NewBool(s.Base.fn("string::is::numeric"))
}

func (s *String[M]) IsSemVer() *Bool[M] {
	return NewBool(s.Base.fn("string::is::semver"))
}

func (s *String[M]) IsURL() *Bool[M] {
	return NewBool(s.Base.fn("string::is::url"))
}

func (s *String[M]) IsUUID() *Bool[M] {
	return NewBool(s.Base.fn("string::is::uuid"))
}

// TODO: string::is::record (needed?)

func (s *String[M]) Base64Decode() *ByteSlice[M] {
	return NewByteSlice(s.Base.fn("encoding::base64::decode"))
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
