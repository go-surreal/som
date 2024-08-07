//go:build embed

package relate

import (
	"context"
)

type Database interface {
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
	Unmarshal(buf []byte, val any) error
}
