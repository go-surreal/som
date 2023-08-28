//go:build embed

package som

import (
	"context"
	"fmt"
	"github.com/marcbinz/som/sdbc"
)

type Database interface {
	Close() error
	Create(ctx context.Context, thing string, data any) (any, error)
	Select(ctx context.Context, what string) (any, error)
	Query(ctx context.Context, statement string, vars map[string]any) (any, error)
	Update(ctx context.Context, thing string, data any) (any, error)
	Delete(ctx context.Context, what string) (any, error)
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

	return &ClientImpl{db: &database{Client: surreal}}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
