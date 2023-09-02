// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/examples/testing/gen/som/conv"
	query "github.com/marcbinz/som/examples/testing/gen/som/query"
	relate "github.com/marcbinz/som/examples/testing/gen/som/relate"
	model "github.com/marcbinz/som/examples/testing/model"
)

type URLExampleRepo interface {
	Query() query.URLExample
	Create(ctx context.Context, user *model.URLExample) error
	CreateWithID(ctx context.Context, id string, user *model.URLExample) error
	Read(ctx context.Context, id string) (*model.URLExample, bool, error)
	Update(ctx context.Context, user *model.URLExample) error
	Delete(ctx context.Context, user *model.URLExample) error
	Relate() *relate.URLExample
}

func (c *ClientImpl) URLExampleRepo() URLExampleRepo {
	return &uRLExample{db: c.db, marshal: c.marshal, unmarshal: c.unmarshal}
}

type uRLExample struct {
	db        Database
	marshal   func(val any) ([]byte, error)
	unmarshal func(buf []byte, val any) error
}

func (n *uRLExample) Query() query.URLExample {
	return query.NewURLExample(n.db, n.unmarshal)
}

func (n *uRLExample) Create(ctx context.Context, uRLExample *model.URLExample) error {
	if uRLExample == nil {
		return errors.New("the passed node must not be nil")
	}
	if uRLExample.ID() != "" {
		return errors.New("given node already has an id")
	}
	key := "url_example"
	data := conv.FromURLExample(*uRLExample)

	raw, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNodes []conv.URLExample
	err = n.unmarshal(raw, &convNodes)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}
	if len(convNodes) < 1 {
		return errors.New("response is empty")
	}
	*uRLExample = conv.ToURLExample(convNodes[0])
	return nil
}

func (n *uRLExample) CreateWithID(ctx context.Context, id string, uRLExample *model.URLExample) error {
	if uRLExample == nil {
		return errors.New("the passed node must not be nil")
	}
	if uRLExample.ID() != "" {
		return errors.New("creating node with preset ID not allowed, use CreateWithID for that")
	}
	key := "url_example:" + "⟨" + id + "⟩"
	data := conv.FromURLExample(*uRLExample)

	res, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNode conv.URLExample
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*uRLExample = conv.ToURLExample(convNode)
	return nil
}

func (n *uRLExample) Read(ctx context.Context, id string) (*model.URLExample, bool, error) {
	res, err := n.db.Select(ctx, "url_example:⟨"+id+"⟩")
	if err != nil {
		return nil, false, fmt.Errorf("could not read entity: %w", err)
	}
	var convNode conv.URLExample
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return nil, false, fmt.Errorf("could not unmarshal entity: %w", err)
	}
	node := conv.ToURLExample(convNode)
	return &node, true, nil
}

func (n *uRLExample) Update(ctx context.Context, uRLExample *model.URLExample) error {
	if uRLExample == nil {
		return errors.New("the passed node must not be nil")
	}
	if uRLExample.ID() == "" {
		return errors.New("cannot update URLExample without existing record ID")
	}
	data := conv.FromURLExample(*uRLExample)

	res, err := n.db.Update(ctx, "url_example:⟨"+uRLExample.ID()+"⟩", data)
	if err != nil {
		return fmt.Errorf("could not update entity: %w", err)
	}
	var convNode conv.URLExample
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*uRLExample = conv.ToURLExample(convNode)
	return nil
}

func (n *uRLExample) Delete(ctx context.Context, uRLExample *model.URLExample) error {
	if uRLExample == nil {
		return errors.New("the passed node must not be nil")
	}
	_, err := n.db.Delete(ctx, "url_example:⟨"+uRLExample.ID()+"⟩")
	if err != nil {
		return fmt.Errorf("could not delete entity: %w", err)
	}
	return nil
}

func (n *uRLExample) Relate() *relate.URLExample {
	return relate.NewURLExample(n.db, n.unmarshal)
}
