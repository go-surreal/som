package sdb

import (
	"errors"
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type database struct {
	*surrealdb.DB
}

func (db *database) Create(thing string, data map[string]any) (any, error) {
	return db.DB.Create(thing, data)
}

func (db *database) Select(what string) (any, error) {
	return db.DB.Select(what)
}

func (db *database) Query(statement string, vars map[string]any) ([]map[string]any, error) {
	fmt.Println(statement)

	raw, err := db.DB.Query(statement, vars)
	if err != nil {
		return nil, err
	}

	fmt.Println(raw)

	if raw == nil {
		return nil, errors.New("database result is nil")
	}

	rawSlice, ok := raw.([]any)
	if !ok {
		return nil, errors.New("database result has invalid format")
	}

	if len(rawSlice) < 1 {
		return nil, errors.New("database result is empty")
	}

	rawMap, ok := raw.([]any)[0].(map[string]any)
	if !ok {
		return nil, errors.New("database result has invalid content")
	}

	status, ok := rawMap["status"]
	if !ok {
		return nil, errors.New("database result does not provide a status")
	}

	if fmt.Sprintf("%s", status) == "ERR" {
		return nil, fmt.Errorf("database returned an error: %s", rawMap["detail"])
	}

	if fmt.Sprintf("%s", status) != "OK" {
		return nil, fmt.Errorf("database returned an unknown status: %s", status)
	}

	rawRows, ok := rawMap["result"].([]any)
	if !ok {
		return nil, errors.New("database result data has invalid format")
	}

	var rows []map[string]any
	for _, row := range rawRows {
		rows = append(rows, row.(map[string]any))
	}

	return rows, nil
}

func (db *database) Update(what string, data map[string]any) (any, error) {
	return db.DB.Update(what, data)
}

func (db *database) Delete(what string) (any, error) {
	return db.DB.Delete(what)
}
