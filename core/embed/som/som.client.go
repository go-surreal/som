//go:build embed

package som

import (
	"context"
	"fmt"
	"github.com/go-surreal/sdbc"
)

type Database interface {
	Create(ctx context.Context, id sdbc.RecordID, data any) ([]byte, error)
	Select(ctx context.Context, id *sdbc.ID) ([]byte, error)
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
	Live(ctx context.Context, statement string, vars map[string]any) (<-chan []byte, error)
	Update(ctx context.Context, id *sdbc.ID, data any) ([]byte, error)
	Delete(ctx context.Context, id *sdbc.ID) ([]byte, error)

	Marshal(val any) ([]byte, error)
	Unmarshal(buf []byte, val any) error
	Close() error
}

type ID *sdbc.ID

type Config struct {
	Host      string
	Secure    bool
	Username  string
	Password  string
	Namespace string
	Database  string
}

type ClientImpl struct {
	db Database
}

func NewClient(ctx context.Context, conf Config, opts ...Option) (*ClientImpl, error) {
	opt := applyOptions(opts)

	surreal, err := sdbc.NewClient(ctx,
		sdbc.Config(conf),
		opt.sdbc...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sdbc client: %v", err)
	}

	return &ClientImpl{
		db: surreal,
	}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
