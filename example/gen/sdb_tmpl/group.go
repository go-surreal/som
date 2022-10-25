package sdb

import (
	"context"
	model "github.com/marcbinz/sdb/db/model"
	query "go.alfnet.dev/service/gampi/gen/sdb/query"
)

var Group group

type group struct{}

func (group) Query() *query.Group {
	return &query.Group{}
}
func (group) Create(ctx context.Context, group *model.Group) error {
	return nil
}
func (group) Update(ctx context.Context, group *model.Group) error {
	return nil
}
func (group) Delete(ctx context.Context, group *model.Group) error {
	return nil
}
