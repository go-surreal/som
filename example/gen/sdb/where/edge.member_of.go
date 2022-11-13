package where

import filter "github.com/marcbinz/sdb/lib/filter"

func newMemberOfIn[T any](key filter.Key) memberOfIn[T] {
	return memberOfIn[T]{memberOf[T]{
		Since: filter.NewTime[T](key.Dot("since")),
		key:   key,
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
		Since: filter.NewTime[T](key.Dot("since")),
		key:   key,
	}}
}

type memberOfOut[T any] struct {
	memberOf[T]
}

func (o memberOfOut[T]) User() user[T] {
	return newUser[T](o.key.Out("user"))
}

type memberOf[T any] struct {
	key   filter.Key
	Since *filter.Time[T]
}

func (n memberOf[T]) Meta() memberOfMeta[T] {
	return newMemberOfMeta[T](n.key.Dot("meta"))
}
