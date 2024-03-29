//go:build embed

package som

import (
	"context"
	"github.com/go-surreal/sdbc"
)

type database struct {
	*sdbc.Client
}

func (db *database) Create(ctx context.Context, thing string, data any) ([]byte, error) {
	return db.Client.Create(ctx, thing, data)
}

func (db *database) Select(ctx context.Context, what string) ([]byte, error) {
	return db.Client.Select(ctx, what)
}

func (db *database) Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error) {
	return db.Client.Query(ctx, statement, vars)
}

func (db *database) Live(ctx context.Context, statement string, vars map[string]any) (<-chan []byte, error) {
	return db.Client.Live(ctx, statement, vars)
}

func (db *database) Update(ctx context.Context, what string, data any) ([]byte, error) {
	return db.Client.Update(ctx, what, data)
}

func (db *database) Delete(ctx context.Context, what string) ([]byte, error) {
	return db.Client.Delete(ctx, what)
}
