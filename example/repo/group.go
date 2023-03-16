package repo

import (
	"github.com/marcbinz/som/example/gen/som"
)

type GroupRepo interface {
	som.GroupRepo
}

type groupRepo struct {
	som.GroupRepo
}

func Group(db som.Client) GroupRepo {
	return &groupRepo{
		GroupRepo: db.GroupRepo(),
	}
}
