package sdb

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type Database interface {
	Close()
	Create(thing string, data map[string]any) (any, error)
	Select(what string) (any, error)
	Query(statement string, vars map[string]any) (any, error)
	Update(thing string, data map[string]any) (any, error)
	Delete(what string) (any, error)
}

type Client struct {
	db Database
}

func NewClient(addr, namespace, database, username, password string) (*Client, error) {
	db, err := surrealdb.New(addr + "/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	_, err = db.Signin(map[string]any{
		"user": username,
		"pass": password,
	})
	if err != nil {
		return nil, err
	}

	_, err = db.Use(namespace, database)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() {
	c.db.Close()
}
