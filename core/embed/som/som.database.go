//go:build embed

package som

import (
	"context"
	"github.com/marcbinz/som/sdbc"
)

type database struct {
	*sdbc.Client
}

func (db *database) Create(ctx context.Context, thing string, data any) (any, error) {
	return db.Client.Create(ctx, 0, thing, data)
}

func (db *database) Select(ctx context.Context, what string) (any, error) {
	return db.Client.Select(ctx, 0, what)
}

func (db *database) Query(ctx context.Context, statement string, vars map[string]any) (any, error) {
	raw, err := db.Client.Query(ctx, 0, statement, vars)
	if err != nil {
		return nil, err
	}

	return raw, err
}

func (db *database) Update(ctx context.Context, what string, data any) (any, error) {
	return db.Client.Update(ctx, 0, what, data)
}

func (db *database) Delete(ctx context.Context, what string) (any, error) {
	return db.Client.Delete(ctx, 0, what)
}
