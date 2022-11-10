package where

import (
	uuid "github.com/google/uuid"
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

var User = newUser[model.User]("")

func newUser[T any](key string) user[T] {
	return user[T]{
		Bool:      filter.NewBool[T](keyed(key, "bool")),
		Bool2:     filter.NewBool[T](keyed(key, "bool_2")),
		CreatedAt: filter.NewTime[T](keyed(key, "created_at")),
		Float32:   filter.NewNumeric[float32, T](keyed(key, "float_32")),
		Float64:   filter.NewNumeric[float64, T](keyed(key, "float_64")),
		ID:        filter.NewID[T](keyed(key, "id"), "user"),
		Int:       filter.NewNumeric[int, T](keyed(key, "int")),
		Int32:     filter.NewNumeric[int32, T](keyed(key, "int_32")),
		Int64:     filter.NewNumeric[int64, T](keyed(key, "int_64")),
		Role:      filter.NewBase[model.Role, T](keyed(key, "role")),
		String:    filter.NewString[T](keyed(key, "string")),
		UUID:      filter.NewBase[uuid.UUID, T](keyed(key, "uuid")),
		UpdatedAt: filter.NewTime[T](keyed(key, "updated_at")),
		key:       key,
	}
}

type user[T any] struct {
	key       string
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
	return newLogin[T](keyed(n.key, "login"))
}
func (n user[T]) Groups() groupSlice[T] {
	key := keyed(n.key, "groups")
	return groupSlice[T]{newGroup[T](key), filter.NewSlice[model.Group, T](key)}
}
func (n user[T]) MainGroup() group[T] {
	return newGroup[T](keyed(n.key, "main_group"))
}
func (n user[T]) Other() *filter.Slice[string, T] {
	return filter.NewSlice[string, T](keyed(n.key, "other"))
}
func (n user[T]) More() *filter.Slice[float32, T] {
	return filter.NewSlice[float32, T](keyed(n.key, "more"))
}
func (n user[T]) Roles() *filter.Slice[model.Role, T] {
	return filter.NewSlice[model.Role, T](keyed(n.key, "roles"))
}
