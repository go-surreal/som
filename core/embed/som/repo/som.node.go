//go:build embed

package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-surreal/sdbc"
)

// N is a placeholder for the node type.
// C is a placeholder for the conversion type.
type repo[N any, C any] struct {
	db Database

	name     string
	convFrom func(*N) *C
	convTo   func(*C) *N
}

func (r *repo[N, C]) create(ctx context.Context, node *N) error {
	data := r.convFrom(node)
	raw, err := r.db.Create(ctx, sdbc.NewULID(r.name), data) // TODO: make ID type configurable
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var conv *C
	err = r.db.Unmarshal(raw, &conv)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}
	*node = *r.convTo(conv)
	return nil
}

func (r *repo[N, C]) createWithID(ctx context.Context, id string, node *N) error {
	data := r.convFrom(node)
	res, err := r.db.Create(ctx, sdbc.MakeID(r.name, id), data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var conv *C
	err = r.db.Unmarshal(res, &conv)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*node = *r.convTo(conv)
	return nil
}

func (r *repo[N, C]) read(ctx context.Context, id *ID) (*N, bool, error) {
	res, err := r.db.Select(ctx, id)
	if err != nil {
		return nil, false, fmt.Errorf("could not read entity: %w", err)
	}
	var conv *C
	err = r.db.Unmarshal(res, &conv)
	if err != nil {
		return nil, false, fmt.Errorf("could not unmarshal entity: %w", err)
	}
	return r.convTo(conv), true, nil
}

func (r *repo[N, C]) update(ctx context.Context, id *ID, node *N) error {
	data := r.convFrom(node)
	res, err := r.db.Update(ctx, id, data)
	if err != nil {
		return fmt.Errorf("could not update entity: %w", err)
	}
	var conv *C
	err = r.db.Unmarshal(res, &conv)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*node = *r.convTo(conv)
	return nil
}

func (r *repo[N, C]) delete(ctx context.Context, id *ID, node *N) error {
	_, err := r.db.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete entity: %w", err)
	}
	return nil
}

func (r *repo[N, C]) refresh(ctx context.Context, id *sdbc.ID, node *N) error {
	read, exists, err := r.read(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to read node: %w", err)
	}

	if !exists {
		return errors.New("given node does not exist")
	}

	*node = *read

	return nil
}
