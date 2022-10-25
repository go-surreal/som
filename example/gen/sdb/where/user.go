package where

import lib "github.com/marcbinz/sdb/lib"

var User = newUser("")

func newUser(origin string) user {
	return user{
		ID:   lib.WhereInt{origin},
		Role: lib.WhereString{origin},
	}
}

type user struct {
	ID   lib.WhereInt
	Role lib.WhereString
}

func (user) Login()        {}
func (user) Groups()       {}
func (user) Other()        {}
func (user) Roles()        {}
func (user) MappedLogin()  {}
func (user) MappedRoles()  {}
func (user) MappedGroups() {}
func (user) OtherMap()     {}
