package where

import (
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

var Group = newGroup("")

func newGroup(key string) group {
	return group{Name: filter.NewString[model.Group](key)}
}

type group struct {
	Name *filter.String[model.Group]
}
