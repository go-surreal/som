package sdb

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	query "github.com/marcbinz/sdb/example/gen/sdb/query"
	model "github.com/marcbinz/sdb/example/model"
)

var Group group

type group struct{}

func (group) Query() *query.Group {
	return query.NewGroup()
}
func (group) Create(ctx context.Context, db *Client, group *model.Group) error {
	if group.ID != "" {
		return errors.New("ID must not be set for a node to be created")
	}
	data := conv.FromGroup(*group)
	raw, err := db.Create("group", data)
	if err != nil {
		return err
	}
	res := conv.ToGroup(raw.([]any)[0].(map[string]any))
	fmt.Println(res)
	return nil
}
func (group) Read(ctx context.Context, db *Client, id string) (*model.Group, error) {
	raw, err := db.db.Select("group:" + id)
	if err != nil {
		return nil, err
	}
	res := conv.ToGroup(raw.(map[string]any))
	return &res, nil
}
func (group) Update(ctx context.Context, group *model.Group) error {
	return nil
}
func (group) Delete(ctx context.Context, group *model.Group) error {
	return nil
}
