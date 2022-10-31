package by

import (
	model "github.com/marcbinz/sdb/example/model"
	sort "github.com/marcbinz/sdb/lib/sort"
)

var Group = newGroup("")

func newGroup(key string) group {
	return group{
		ID:   sort.NewSort[model.Group](key),
		Name: sort.NewString[model.Group](key),
	}
}

type group struct {
	ID   *sort.Sort[model.Group]
	Name *sort.String[model.Group]
}

func (group) Random() *sort.Of[model.Group] {
	return nil
}
