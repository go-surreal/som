// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package relate

import (
	"context"
)

type Database interface {
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
}
