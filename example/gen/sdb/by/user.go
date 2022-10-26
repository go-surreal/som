package by

import (
	model "github.com/marcbinz/sdb/example/model"
	sort "github.com/marcbinz/sdb/lib/sort"
)

var User = newUser("")

func newUser(key string) user {
	return user{
		ID:   sort.NewSort[model.User](key),
		Text: sort.NewString[model.User](key),
	}
}

type user struct {
	ID   *sort.Sort[model.User]
	Text *sort.String[model.User]
}
