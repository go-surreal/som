//go:build embed

package som

import (
	"context"
	"fmt"
)

// N is a placeholder for the node type.
// C is a placeholder for the conversion type.
type repo[N any, C any] struct {
	db        Database
	marshal   func(val any) ([]byte, error)
	unmarshal func(buf []byte, val any) error

	name     string
	convFrom func(*N) *C
	convTo   func(*C) *N
}

func (r *repo[N, C]) create(ctx context.Context, node *N) error {
	key := r.name + ":ulid()"
	data := r.convFrom(node)
	raw, err := r.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var conv *C
	err = r.unmarshal(raw, &conv)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}
	*node = *r.convTo(conv)
	return nil
}

func (r *repo[N, C]) createWithID(ctx context.Context, id string, node *N) error {
	key := r.name + ":" + "⟨" + id + "⟩"
	data := r.convFrom(node)
	res, err := r.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var conv *C
	err = r.unmarshal(res, &conv)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*node = *r.convTo(conv)
	return nil
}

func (r *repo[N, C]) read(ctx context.Context, id string) (*N, bool, error) {
	res, err := r.db.Select(ctx, r.name+":⟨"+id+"⟩")
	if err != nil {
		return nil, false, fmt.Errorf("could not read entity: %w", err)
	}
	var conv *C
	err = r.unmarshal(res, &conv)
	if err != nil {
		return nil, false, fmt.Errorf("could not unmarshal entity: %w", err)
	}
	return r.convTo(conv), true, nil
}

func (r *repo[N, C]) update(ctx context.Context, id string, node *N) error {
	data := r.convFrom(node)
	res, err := r.db.Update(ctx, r.name+":⟨"+id+"⟩", data)
	if err != nil {
		return fmt.Errorf("could not update entity: %w", err)
	}
	var conv *C
	err = r.unmarshal(res, &conv)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*node = *r.convTo(conv)
	return nil
}

func (r *repo[N, C]) delete(ctx context.Context, id string, node *N) error {
	_, err := r.db.Delete(ctx, r.name+":⟨"+id+"⟩")
	if err != nil {
		return fmt.Errorf("could not delete entity: %w", err)
	}
	return nil
}
