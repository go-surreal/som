// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package som

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type Database interface {
	Close()
	Create(thing string, data any) (any, error)
	Select(what string) (any, error)
	Query(statement string, vars any) (any, error)
	Update(thing string, data any) (any, error)
	Delete(what string) (any, error)
}

type Config struct {
	Address string
	Username string
	Password string
	Namespace string
	Database string
}

type Client interface {
	User() UserRepo
}

type ClientImpl struct {
	db Database
}

func NewClient(conf Config) (*ClientImpl, error) {
	surreal, err := surrealdb.New(conf.Address + "/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	_, err = surreal.Signin(map[string]any{
		"user": conf.Username,
		"pass": conf.Password,
	})
	if err != nil {
		return nil, err
	}

	_, err = surreal.Use(conf.Namespace, conf.Database)
	if err != nil {
		return nil, err
	}

	return &ClientImpl{db: &database{DB: surreal}}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
