// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/go-surreal/som/examples/testing/gen/som/conv"
	query "github.com/go-surreal/som/examples/testing/gen/som/query"
	relate "github.com/go-surreal/som/examples/testing/gen/som/relate"
	model "github.com/go-surreal/som/examples/testing/model"
)

type FieldsLikeDBResponseRepo interface {
	Query() query.NodeFieldsLikeDBResponse
	Create(ctx context.Context, user *model.FieldsLikeDBResponse) error
	CreateWithID(ctx context.Context, id string, user *model.FieldsLikeDBResponse) error
	Read(ctx context.Context, id string) (*model.FieldsLikeDBResponse, bool, error)
	Update(ctx context.Context, user *model.FieldsLikeDBResponse) error
	Delete(ctx context.Context, user *model.FieldsLikeDBResponse) error
	Relate() *relate.FieldsLikeDBResponse
}

func (c *ClientImpl) FieldsLikeDBResponseRepo() FieldsLikeDBResponseRepo {
	return &fieldsLikeDBResponse{db: c.db, marshal: c.marshal, unmarshal: c.unmarshal}
}

type fieldsLikeDBResponse struct {
	db        Database
	marshal   func(val any) ([]byte, error)
	unmarshal func(buf []byte, val any) error
}

func (n *fieldsLikeDBResponse) Query() query.NodeFieldsLikeDBResponse {
	return query.NewFieldsLikeDBResponse(n.db, n.unmarshal)
}

func (n *fieldsLikeDBResponse) Create(ctx context.Context, fieldsLikeDBResponse *model.FieldsLikeDBResponse) error {
	if fieldsLikeDBResponse == nil {
		return errors.New("the passed node must not be nil")
	}
	if fieldsLikeDBResponse.ID() != "" {
		return errors.New("given node already has an id")
	}
	key := "fields_like_db_response"
	data := conv.FromFieldsLikeDBResponse(fieldsLikeDBResponse)

	raw, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNodes []*conv.FieldsLikeDBResponse
	err = n.unmarshal(raw, &convNodes)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}
	if len(convNodes) < 1 {
		return errors.New("response is empty")
	}
	*fieldsLikeDBResponse = *conv.ToFieldsLikeDBResponse(convNodes[0])
	return nil
}

func (n *fieldsLikeDBResponse) CreateWithID(ctx context.Context, id string, fieldsLikeDBResponse *model.FieldsLikeDBResponse) error {
	if fieldsLikeDBResponse == nil {
		return errors.New("the passed node must not be nil")
	}
	if fieldsLikeDBResponse.ID() != "" {
		return errors.New("creating node with preset ID not allowed, use CreateWithID for that")
	}
	key := "fields_like_db_response:" + "⟨" + id + "⟩"
	data := conv.FromFieldsLikeDBResponse(fieldsLikeDBResponse)

	res, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNode *conv.FieldsLikeDBResponse
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*fieldsLikeDBResponse = *conv.ToFieldsLikeDBResponse(convNode)
	return nil
}

func (n *fieldsLikeDBResponse) Read(ctx context.Context, id string) (*model.FieldsLikeDBResponse, bool, error) {
	res, err := n.db.Select(ctx, "fields_like_db_response:⟨"+id+"⟩")
	if err != nil {
		return nil, false, fmt.Errorf("could not read entity: %w", err)
	}
	var convNode *conv.FieldsLikeDBResponse
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return nil, false, fmt.Errorf("could not unmarshal entity: %w", err)
	}
	return conv.ToFieldsLikeDBResponse(convNode), true, nil
}

func (n *fieldsLikeDBResponse) Update(ctx context.Context, fieldsLikeDBResponse *model.FieldsLikeDBResponse) error {
	if fieldsLikeDBResponse == nil {
		return errors.New("the passed node must not be nil")
	}
	if fieldsLikeDBResponse.ID() == "" {
		return errors.New("cannot update FieldsLikeDBResponse without existing record ID")
	}
	data := conv.FromFieldsLikeDBResponse(fieldsLikeDBResponse)

	res, err := n.db.Update(ctx, "fields_like_db_response:⟨"+fieldsLikeDBResponse.ID()+"⟩", data)
	if err != nil {
		return fmt.Errorf("could not update entity: %w", err)
	}
	var convNode *conv.FieldsLikeDBResponse
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*fieldsLikeDBResponse = *conv.ToFieldsLikeDBResponse(convNode)
	return nil
}

func (n *fieldsLikeDBResponse) Delete(ctx context.Context, fieldsLikeDBResponse *model.FieldsLikeDBResponse) error {
	if fieldsLikeDBResponse == nil {
		return errors.New("the passed node must not be nil")
	}
	_, err := n.db.Delete(ctx, "fields_like_db_response:⟨"+fieldsLikeDBResponse.ID()+"⟩")
	if err != nil {
		return fmt.Errorf("could not delete entity: %w", err)
	}
	return nil
}

func (n *fieldsLikeDBResponse) Relate() *relate.FieldsLikeDBResponse {
	return relate.NewFieldsLikeDBResponse(n.db, n.unmarshal)
}
