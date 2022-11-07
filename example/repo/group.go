package repo

import (
	"context"

	"github.com/marcbinz/sdb/example/gen/sdb"
	"github.com/marcbinz/sdb/example/model"
)

type GroupRepo interface {
	Create(ctx context.Context, user *model.Group) error
}

type groupRepo struct {
	db *sdb.Client
}

func Group(db *sdb.Client) GroupRepo {
	return &groupRepo{db: db}
}

func (repo *groupRepo) Create(ctx context.Context, user *model.Group) error {
	return repo.db.Group().Create(ctx, user)
}
