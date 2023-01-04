// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package where

import lib "github.com/marcbinz/som/lib"

func newMemberOfIn[T any](key lib.Key) memberOfIn[T] {
	return memberOfIn[T]{memberOf[T]{
		CreatedAt: lib.NewTime[T](key.Field("created_at")),
		ID:        lib.NewID[T](key.Field("id"), "member_of"),
		UpdatedAt: lib.NewTime[T](key.Field("updated_at")),
		key:       key,
	}}
}

type memberOfIn[T any] struct {
	memberOf[T]
}

func (i memberOfIn[T]) Group() group[T] {
	return newGroup[T](i.key.EdgeIn("group", nil))
}
func newMemberOfOut[T any](key lib.Key) memberOfOut[T] {
	return memberOfOut[T]{memberOf[T]{
		CreatedAt: lib.NewTime[T](key.Field("created_at")),
		ID:        lib.NewID[T](key.Field("id"), "member_of"),
		UpdatedAt: lib.NewTime[T](key.Field("updated_at")),
		key:       key,
	}}
}

type memberOfOut[T any] struct {
	memberOf[T]
}

func (o memberOfOut[T]) User() user[T] {
	return newUser[T](o.key.EdgeOut("user", nil))
}

type memberOf[T any] struct {
	key       lib.Key
	ID        *lib.ID[T]
	CreatedAt *lib.Time[T]
	UpdatedAt *lib.Time[T]
}

func (n memberOf[T]) Meta() memberOfMeta[T] {
	return newMemberOfMeta[T](n.key.Field("meta"))
}
