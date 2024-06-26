// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package som

import (
	"context"
	"fmt"
	"github.com/go-surreal/sdbc"
)

type Database interface {
	Close() error
	Create(ctx context.Context, thing string, data any) ([]byte, error)
	Select(ctx context.Context, what string) ([]byte, error)
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
	Live(ctx context.Context, statement string, vars map[string]any) (<-chan []byte, error)
	Update(ctx context.Context, thing string, data any) ([]byte, error)
	Delete(ctx context.Context, what string) ([]byte, error)
}

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

	marshal   func(val any) ([]byte, error)
	unmarshal func(buf []byte, val any) error
}

func NewClient(ctx context.Context, conf Config, opts ...Option) (*ClientImpl, error) {
	opt := applyOptions(opts)

	surreal, err := sdbc.NewClient(ctx,
		sdbc.Config(conf),
		opt.sdbc...,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create sdbc client: %v", err)
	}

	return &ClientImpl{
		db: &database{Client: surreal},

		marshal:   opt.jsonMarshal,
		unmarshal: opt.jsonUnmarshal,
	}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
