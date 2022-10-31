package sdb

import surrealdbgo "github.com/surrealdb/surrealdb.go"

type Client struct {
	db *surrealdbgo.DB
}

func NewClient(db *surrealdbgo.DB) *Client {
	return &Client{db: db}
}
func (c *Client) Create(node string, data map[string]any) (any, error) {
	return c.db.Create(node, data)
}
