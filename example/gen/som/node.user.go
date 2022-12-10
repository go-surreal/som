package som

import (
	"context"
	"errors"
	conv "github.com/marcbinz/som/example/gen/som/conv"
	query "github.com/marcbinz/som/example/gen/som/query"
	relate "github.com/marcbinz/som/example/gen/som/relate"
	model "github.com/marcbinz/som/example/model"
	surrealdbgo "github.com/surrealdb/surrealdb.go"
)

func (c *Client) User() *user {
	return &user{client: c}
}

type user struct {
	client *Client
}

func (n *user) Query() *query.User {
	return query.NewUser(n.client.db)
}
func (n *user) Create(ctx context.Context, user *model.User) error {
	if user.ID != "" {
		return errors.New("ID must not be set for a node to be created")
	}
	data := conv.FromUser(user)
	raw, err := n.client.db.Create("user", data)
	if err != nil {
		return err
	}
	var convNode conv.User
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return err
	}
	*user = *conv.ToUser(&convNode)
	return nil
}
func (n *user) Read(ctx context.Context, id string) (*model.User, bool, error) {
	raw, err := n.client.db.Select("user:" + id)
	if err != nil {
		if errors.As(err, &surrealdbgo.PermissionError{}) {
			return nil, false, nil
		}
		return nil, false, err
	}
	var convNode *conv.User
	err = surrealdbgo.Unmarshal([]any{raw}, &convNode)
	if err != nil {
		return nil, false, err
	}
	return conv.ToUser(convNode), true, nil
}
func (n *user) Update(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		return errors.New("cannot update User without existing record ID")
	}
	data := conv.FromUser(user)
	raw, err := n.client.db.Update("user:"+user.ID, data)
	if err != nil {
		return err
	}
	var convNode conv.User
	err = surrealdbgo.Unmarshal([]any{raw}, &convNode)
	if err != nil {
		return err
	}
	*user = *conv.ToUser(&convNode)
	return nil
}
func (n *user) Delete(ctx context.Context, user *model.User) error {
	_, err := n.client.db.Delete("user:" + user.ID)
	if err != nil {
		return err
	}
	return nil
}
func (n *user) Relate() *relate.User {
	return relate.NewUser(n.client.db)
}
