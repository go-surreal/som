package lib

import (
	"github.com/go-surreal/som"
)

// Password is a filter builder for password values.
// M is the model this filter is for.
type Password[M any] struct {
	*Base[M, som.Password]
}

func NewPassword[M any](key Key[M]) *Password[M] {
	return &Password[M]{
		Base: NewBase[M, som.Password](key),
	}
}

func (e *Email[M]) CompareArgon2(val string) *Bool[M] {
	return NewBool(e.key.fn("crypto::argon2::compare", val))
}

func (e *Email[M]) CompareBcrypt(val string) *Bool[M] {
	return NewBool(e.key.fn("crypto::bcrypt::compare", val))
}

func (e *Email[M]) ComparePbkdf2(val string) *Bool[M] {
	return NewBool(e.key.fn("crypto::pbkdf2::compare", val))
}

func (e *Email[M]) CompareScrypt(val string) *Bool[M] {
	return NewBool(e.key.fn("crypto::scrypt::compare", val))
}
