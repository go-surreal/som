package sdb

import (
	// "errors"
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type database struct {
	*surrealdb.DB
}

func (db *database) Create(thing string, data any) (any, error) {
	return db.DB.Create(thing, data)
}

func (db *database) Select(what string) (any, error) {
	return db.DB.Select(what)
}

func (db *database) Query(statement string, vars map[string]any) (any, error) {
	fmt.Println(statement)

	raw, err := db.DB.Query(statement, vars)
	if err != nil {
		return nil, err
	}

	fmt.Println(raw)
	
	return raw, err
}

func (db *database) Update(what string, data map[string]any) (any, error) {
	return db.DB.Update(what, data)
}

func (db *database) Delete(what string) (any, error) {
	return db.DB.Delete(what)
}
