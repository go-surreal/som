package repo

import (
	"github.com/go-surreal/som/tests/basic/gen/som/repo"
)

type GroupRepo interface {
	repo.GroupRepo
}

type groupRepo struct {
	repo.GroupRepo
}

func Group(db repo.Client) GroupRepo {
	return &groupRepo{
		GroupRepo: db.GroupRepo(),
	}
}
