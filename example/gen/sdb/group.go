package sdb

import (
	"context"
	query "github.com/marcbinz/sdb/example/gen/sdb/query"
	model "github.com/marcbinz/sdb/example/model"
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
