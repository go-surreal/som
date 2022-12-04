package where

import (
	uuid "github.com/google/uuid"
	model "github.com/marcbinz/som/example/model"
	filter "github.com/marcbinz/som/lib/filter"
)

var User = newUser[model.User](filter.NewKey())

func newUser[T any](key filter.Key) user[T] {
	return user[T]{
		Bool:      filter.NewBool[T](key.Dot("bool")),
		Bool2:     filter.NewBool[T](key.Dot("bool_2")),
		CreatedAt: filter.NewTime[T](key.Dot("created_at")),
		Float32:   filter.NewNumeric[float32, T](key.Dot("float_32")),
		Float64:   filter.NewNumeric[float64, T](key.Dot("float_64")),
		ID:        filter.NewID[T](key.Dot("id"), "user"),
		Int:       filter.NewNumeric[int, T](key.Dot("int")),
		Int32:     filter.NewNumeric[int32, T](key.Dot("int_32")),
		Int64:     filter.NewNumeric[int64, T](key.Dot("int_64")),
		Role:      filter.NewBase[model.Role, T](key.Dot("role")),
		String:    filter.NewString[T](key.Dot("string")),
		UUID:      filter.NewBase[uuid.UUID, T](key.Dot("uuid")),
		UpdatedAt: filter.NewTime[T](key.Dot("updated_at")),
		key:       key,
	}
}

type user[T any] struct {
	key       filter.Key
	ID        *filter.ID[T]
	CreatedAt *filter.Time[T]
	UpdatedAt *filter.Time[T]
	String    *filter.String[T]
	Int       *filter.Numeric[int, T]
	Int32     *filter.Numeric[int32, T]
	Int64     *filter.Numeric[int64, T]
	Float32   *filter.Numeric[float32, T]
	Float64   *filter.Numeric[float64, T]
	Bool      *filter.Bool[T]
	Bool2     *filter.Bool[T]
	UUID      *filter.Base[uuid.UUID, T]
	Role      *filter.Base[model.Role, T]
}
type userSlice[T any] struct {
	user[T]
	*filter.Slice[model.User, T]
}

func (n user[T]) Login() login[T] {
	return newLogin[T](n.key.Dot("login"))
}
func (n user[T]) Groups() groupSlice[T] {
	key := n.key.Dot("groups")
	return groupSlice[T]{newGroup[T](key), filter.NewSlice[model.Group, T](key)}
}
func (n user[T]) MainGroup() group[T] {
	return newGroup[T](n.key.Dot("main_group"))
}
func (n user[T]) Other() *filter.Slice[string, T] {
	return filter.NewSlice[string, T](n.key.Dot("other"))
}
func (n user[T]) More() *filter.Slice[float32, T] {
	return filter.NewSlice[float32, T](n.key.Dot("more"))
}
func (n user[T]) Roles() *filter.Slice[model.Role, T] {
	return filter.NewSlice[model.Role, T](n.key.Dot("roles"))
}
func (n user[T]) MyGroups() memberOfIn[T] {
	return newMemberOfIn[T](n.key.In("member_of"))
}