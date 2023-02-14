package repo

import (
	"context"
	"github.com/marcbinz/som/examples/basic/gen/som"
	"github.com/marcbinz/som/examples/basic/model"
)

type GroupRepo interface {
	Create(ctx context.Context, user *model.Group) error
}

type groupRepo struct {
	db *som.Client
}

func Group(db *som.Client) GroupRepo {
	return &groupRepo{db: db}
}

func (repo *groupRepo) Create(ctx context.Context, user *model.Group) error {
	return repo.db.Group().Create(ctx, user)
}
