package where

import (
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

var User = newUser[model.User]("")

func newUser[T any](key string) user[T] {
	return user[T]{
		Bool:      filter.NewBool[T](key),
		CreatedAt: filter.NewTime[T](key),
		Float32:   filter.NewNumeric[float32, T](key),
		Float64:   filter.NewNumeric[float64, T](key),
		ID:        filter.NewBase[string, T](key),
		Int:       filter.NewNumeric[int, T](key),
		Int32:     filter.NewNumeric[int32, T](key),
		Int64:     filter.NewNumeric[int64, T](key),
		Role:      filter.NewBase[model.Role, T](key),
		String:    filter.NewString[T](key),
		UpdatedAt: filter.NewTime[T](key),
	}
}

type user[T any] struct {
	ID        *filter.Base[string, T]
	CreatedAt *filter.Time[T]
	UpdatedAt *filter.Time[T]
	String    *filter.String[T]
	Int       *filter.Numeric[int, T]
	Int32     *filter.Numeric[int32, T]
	Int64     *filter.Numeric[int64, T]
	Float32   *filter.Numeric[float32, T]
	Float64   *filter.Numeric[float64, T]
	Bool      *filter.Bool[T]
	Role      *filter.Base[model.Role, T]
}

func (user[T]) Login()  {}
func (user[T]) Groups() {}
func (user[T]) MainGroup() group[T] {
	return newGroup[T]("main_group")
}
func (user[T]) Other()        {}
func (user[T]) Roles()        {}
func (user[T]) MappedLogin()  {}
func (user[T]) MappedRoles()  {}
func (user[T]) MappedGroups() {}
func (user[T]) OtherMap()     {}
