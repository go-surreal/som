package with

import model "github.com/marcbinz/som/example/model"

var Group = group[model.Group]("")

type group[T any] string

func (n group[T]) fetch(T) {}