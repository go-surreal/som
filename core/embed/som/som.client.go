//go:build embed

package som

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/marcbinz/som/sdbc"
)

type Database interface {
	Close() error
	Create(ctx context.Context, thing string, data any) ([]byte, error)
	Select(ctx context.Context, what string) ([]byte, error)
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
	Update(ctx context.Context, thing string, data any) ([]byte, error)
	Delete(ctx context.Context, what string) ([]byte, error)
}

type Config struct {
	Address   string
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

func NewClient(ctx context.Context, conf Config) (*ClientImpl, error) {
	url := conf.Address + "/rpc"

	surreal, err := sdbc.NewClient(ctx, sdbc.Config{
		Address:   url,
		Username:  conf.Username,
		Password:  conf.Password,
		Namespace: conf.Namespace,
		Database:  conf.Database,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create sdbc client: %v", err)
	}

	return &ClientImpl{
		db: &database{Client: surreal},

		marshal:   sonic.ConfigFastest.Marshal,
		unmarshal: sonic.ConfigFastest.Unmarshal,
	}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
