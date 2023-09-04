// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/examples/movie/gen/som/conv"
	query "github.com/marcbinz/som/examples/movie/gen/som/query"
	relate "github.com/marcbinz/som/examples/movie/gen/som/relate"
	model "github.com/marcbinz/som/examples/movie/model"
)

type MovieRepo interface {
	Query() query.Movie
	Create(ctx context.Context, user *model.Movie) error
	CreateWithID(ctx context.Context, id string, user *model.Movie) error
	Read(ctx context.Context, id string) (*model.Movie, bool, error)
	Update(ctx context.Context, user *model.Movie) error
	Delete(ctx context.Context, user *model.Movie) error
	Relate() *relate.Movie
}

func (c *ClientImpl) MovieRepo() MovieRepo {
	return &movie{db: c.db, marshal: c.marshal, unmarshal: c.unmarshal}
}

type movie struct {
	db        Database
	marshal   func(val any) ([]byte, error)
	unmarshal func(buf []byte, val any) error
}

func (n *movie) Query() query.Movie {
	return query.NewMovie(n.db, n.unmarshal)
}

func (n *movie) Create(ctx context.Context, movie *model.Movie) error {
	if movie == nil {
		return errors.New("the passed node must not be nil")
	}
	if movie.ID() != "" {
		return errors.New("given node already has an id")
	}
	key := "movie"
	data := conv.FromMovie(movie)

	raw, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNodes []*conv.Movie
	err = n.unmarshal(raw, &convNodes)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}
	if len(convNodes) < 1 {
		return errors.New("response is empty")
	}
	*movie = *conv.ToMovie(convNodes[0])
	return nil
}

func (n *movie) CreateWithID(ctx context.Context, id string, movie *model.Movie) error {
	if movie == nil {
		return errors.New("the passed node must not be nil")
	}
	if movie.ID() != "" {
		return errors.New("creating node with preset ID not allowed, use CreateWithID for that")
	}
	key := "movie:" + "⟨" + id + "⟩"
	data := conv.FromMovie(movie)

	res, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNode *conv.Movie
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*movie = *conv.ToMovie(convNode)
	return nil
}

func (n *movie) Read(ctx context.Context, id string) (*model.Movie, bool, error) {
	res, err := n.db.Select(ctx, "movie:⟨"+id+"⟩")
	if err != nil {
		return nil, false, fmt.Errorf("could not read entity: %w", err)
	}
	var convNode *conv.Movie
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return nil, false, fmt.Errorf("could not unmarshal entity: %w", err)
	}
	return conv.ToMovie(convNode), true, nil
}

func (n *movie) Update(ctx context.Context, movie *model.Movie) error {
	if movie == nil {
		return errors.New("the passed node must not be nil")
	}
	if movie.ID() == "" {
		return errors.New("cannot update Movie without existing record ID")
	}
	data := conv.FromMovie(movie)

	res, err := n.db.Update(ctx, "movie:⟨"+movie.ID()+"⟩", data)
	if err != nil {
		return fmt.Errorf("could not update entity: %w", err)
	}
	var convNode *conv.Movie
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*movie = *conv.ToMovie(convNode)
	return nil
}

func (n *movie) Delete(ctx context.Context, movie *model.Movie) error {
	if movie == nil {
		return errors.New("the passed node must not be nil")
	}
	_, err := n.db.Delete(ctx, "movie:⟨"+movie.ID()+"⟩")
	if err != nil {
		return fmt.Errorf("could not delete entity: %w", err)
	}
	return nil
}

func (n *movie) Relate() *relate.Movie {
	return relate.NewMovie(n.db, n.unmarshal)
}
