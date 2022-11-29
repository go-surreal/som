package where

import filter "github.com/marcbinz/sdb/lib/filter"

func newMemberOfIn[T any](key filter.Key) memberOfIn[T] {
	return memberOfIn[T]{memberOf[T]{
		CreatedAt: filter.NewTime[T](key.Dot("created_at")),
		ID:        filter.NewID[T](key.Dot("id"), "member_of"),
		UpdatedAt: filter.NewTime[T](key.Dot("updated_at")),
		key:       key,
	}}
}

type memberOfIn[T any] struct {
	memberOf[T]
}

func (i memberOfIn[T]) Group() group[T] {
	return newGroup[T](i.key.In("group"))
}
func newMemberOfOut[T any](key filter.Key) memberOfOut[T] {
	return memberOfOut[T]{memberOf[T]{
		CreatedAt: filter.NewTime[T](key.Dot("created_at")),
		ID:        filter.NewID[T](key.Dot("id"), "member_of"),
		UpdatedAt: filter.NewTime[T](key.Dot("updated_at")),
		key:       key,
	}}
}

type memberOfOut[T any] struct {
	memberOf[T]
}

func (o memberOfOut[T]) User() user[T] {
	return newUser[T](o.key.Out("user"))
}

type memberOf[T any] struct {
	key       filter.Key
	ID        *filter.ID[T]
	CreatedAt *filter.Time[T]
	UpdatedAt *filter.Time[T]
}

func (n memberOf[T]) Meta() memberOfMeta[T] {
	return newMemberOfMeta[T](n.key.Dot("meta"))
}