package where

import (
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

var User = newUser("")

func newUser(key string) user {
	return user{
		ID:   filter.NewBase[string, model.User](key),
		Role: filter.NewBase[model.Role, model.User](key),
		Text: filter.NewString[model.User](key),
	}
}

type user struct {
	ID   *filter.Base[string, model.User]
	Text *filter.String[model.User]
	Role *filter.Base[model.Role, model.User]
}

func (user) Login()        {}
func (user) Groups()       {}
func (user) Other()        {}
func (user) Roles()        {}
func (user) MappedLogin()  {}
func (user) MappedRoles()  {}
func (user) MappedGroups() {}
func (user) OtherMap()     {}
