package sdb

import surrealdbgo "github.com/surrealdb/surrealdb.go"

type Client struct {
	db *surrealdbgo.DB
}
