package where

import (
	"github.com/marcbinz/sdb/lib"
)

var User = newUser("user")

type user struct {
	ID           lib.WhereID
	Login        struct{}
	Role         struct{}
	Groups       struct{}
	Other        struct{}
	Roles        struct{}
	MappedLogin  struct{}
	MappedRoles  struct{}
	MappedGroups struct{}
	OtherMap     struct{}
}

func newUser(origin string) user {
	return user{
		ID: lib.WhereID{Origin: origin},
	}
}
