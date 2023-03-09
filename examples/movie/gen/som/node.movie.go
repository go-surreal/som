// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	conv "github.com/marcbinz/som/examples/movie/gen/som/conv"
	query "github.com/marcbinz/som/examples/movie/gen/som/query"
	relate "github.com/marcbinz/som/examples/movie/gen/som/relate"
	model "github.com/marcbinz/som/examples/movie/model"
	surrealdbgo "github.com/surrealdb/surrealdb.go"
)

func (c *Client) Movie() *movie {
	return &movie{client: c}
}

type movie struct {
	client *Client
}

func (n *movie) Query() *query.Movie {
	return query.NewMovie(n.client.db)
}

func (n *movie) Create(ctx context.Context, movie *model.Movie) error {
	if movie == nil {
		return errors.New("the passed node must not be nil")
	}
	if movie.ID() != "" {
		return errors.New("creating node with preset ID not allowed, use CreateWithID for that")
	}
	key := "movie"
	data := conv.FromMovie(*movie)

	raw, err := n.client.db.Create(key, data)
	if err != nil {
		return err
	}
	if _, ok := raw.([]any); !ok {
		raw = []any{raw} // temporary fix
	}
	var convNode conv.Movie
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return err
	}
	*movie = conv.ToMovie(convNode)
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
	data := conv.FromMovie(*movie)

	raw, err := n.client.db.Create(key, data)
	if err != nil {
		return err
	}
	if _, ok := raw.([]any); !ok {
		raw = []any{raw} // temporary fix
	}
	var convNode conv.Movie
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return err
	}
	*movie = conv.ToMovie(convNode)
	return nil
}

func (n *movie) Read(ctx context.Context, id string) (*model.Movie, bool, error) {
	raw, err := n.client.db.Select("movie:⟨" + id + "⟩")
	if err != nil {
		if errors.As(err, &surrealdbgo.PermissionError{}) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if _, ok := raw.([]any); !ok {
		raw = []any{raw} // temporary fix
	}
	var convNode conv.Movie
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return nil, false, err
	}
	node := conv.ToMovie(convNode)
	return &node, true, nil
}

func (n *movie) Update(ctx context.Context, movie *model.Movie) error {
	if movie == nil {
		return errors.New("the passed node must not be nil")
	}
	if movie.ID() == "" {
		return errors.New("cannot update Movie without existing record ID")
	}
	data := conv.FromMovie(*movie)

	raw, err := n.client.db.Update("movie:⟨"+movie.ID()+"⟩", data)
	if err != nil {
		return err
	}
	var convNode conv.Movie
	err = surrealdbgo.Unmarshal([]any{raw}, &convNode)
	if err != nil {
		return err
	}
	*movie = conv.ToMovie(convNode)
	return nil
}

func (n *movie) Delete(ctx context.Context, movie *model.Movie) error {
	if movie == nil {
		return errors.New("the passed node must not be nil")
	}
	_, err := n.client.db.Delete("movie:⟨" + movie.ID() + "⟩")
	if err != nil {
		return err
	}
	return nil
}

func (n *movie) Relate() *relate.Movie {
	return relate.NewMovie(n.client.db)
}