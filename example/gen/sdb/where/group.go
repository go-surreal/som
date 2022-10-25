package where

import lib "github.com/marcbinz/sdb/lib"

var Group = newGroup("")

func newGroup(origin string) group {
	return group{Name: lib.WhereString{origin}}
}

type group struct {
	Name lib.WhereString
}
