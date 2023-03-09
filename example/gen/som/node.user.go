// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	conv "github.com/marcbinz/som/example/gen/som/conv"
	query "github.com/marcbinz/som/example/gen/som/query"
	relate "github.com/marcbinz/som/example/gen/som/relate"
	model "github.com/marcbinz/som/example/model"
	surrealdbgo "github.com/surrealdb/surrealdb.go"
	"time"
)

func (c *ClientImpl) User() UserRepo {
	return &user{client: c}
}

type UserRepo interface {
	Query() query.UserQuery
	Create(ctx context.Context, user *model.User) error
	CreateWithID(ctx context.Context, id string, user *model.User) error
	Read(ctx context.Context, id string) (*model.User, bool, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
	Relate() *relate.User
}

type user struct {
	client *ClientImpl
}

func (n *user) Query() query.UserQuery {
	return query.NewUser(n.client.db)
}

func (n *user) Create(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("the passed node must not be nil")
	}
	if user.ID() != "" {
		return errors.New("creating node with preset ID not allowed, use CreateWithID for that")
	}
	key := "user"
	data := conv.FromUser(*user)
	data.CreatedAt = time.Now()
	data.UpdatedAt = data.CreatedAt
	raw, err := n.client.db.Create(key, data)
	if err != nil {
		return err
	}
	if _, ok := raw.([]any); !ok {
		raw = []any{raw} // temporary fix
	}
	var convNode conv.User
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return err
	}
	*user = conv.ToUser(convNode)
	return nil
}

func (n *user) CreateWithID(ctx context.Context, id string, user *model.User) error {
	if user == nil {
		return errors.New("the passed node must not be nil")
	}
	if user.ID() != "" {
		return errors.New("creating node with preset ID not allowed, use CreateWithID for that")
	}
	key := "user:" + "⟨" + id + "⟩"
	data := conv.FromUser(*user)
	data.CreatedAt = time.Now()
	data.UpdatedAt = data.CreatedAt
	raw, err := n.client.db.Create(key, data)
	if err != nil {
		return err
	}
	if _, ok := raw.([]any); !ok {
		raw = []any{raw} // temporary fix
	}
	var convNode conv.User
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return err
	}
	*user = conv.ToUser(convNode)
	return nil
}

func (n *user) Read(ctx context.Context, id string) (*model.User, bool, error) {
	raw, err := n.client.db.Select("user:⟨" + id + "⟩")
	if err != nil {
		if errors.As(err, &surrealdbgo.PermissionError{}) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if _, ok := raw.([]any); !ok {
		raw = []any{raw} // temporary fix
	}
	var convNode conv.User
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return nil, false, err
	}
	node := conv.ToUser(convNode)
	return &node, true, nil
}

func (n *user) Update(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("the passed node must not be nil")
	}
	if user.ID() == "" {
		return errors.New("cannot update User without existing record ID")
	}
	data := conv.FromUser(*user)
	data.UpdatedAt = time.Now()
	raw, err := n.client.db.Update("user:⟨"+user.ID()+"⟩", data)
	if err != nil {
		return err
	}
	var convNode conv.User
	err = surrealdbgo.Unmarshal([]any{raw}, &convNode)
	if err != nil {
		return err
	}
	*user = conv.ToUser(convNode)
	return nil
}

func (n *user) Delete(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("the passed node must not be nil")
	}
	_, err := n.client.db.Delete("user:⟨" + user.ID() + "⟩")
	if err != nil {
		return err
	}
	return nil
}

func (n *user) Relate() *relate.User {
	return relate.NewUser(n.client.db)
}
