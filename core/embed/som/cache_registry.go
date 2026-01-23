//go:build embed

package som

import (
	"reflect"
	"sync"
)

var (
	tableRegistry = make(map[reflect.Type]string)
	registryMu    sync.RWMutex
)

// RegisterTable registers a model type with its corresponding database table name.
// This is typically called from init() functions in generated code.
func RegisterTable[T any](tableName string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	var zero T
	tableRegistry[reflect.TypeOf(zero)] = tableName
}

// getTableName returns the database table name for a model type.
// Returns empty string if the type is not registered.
func getTableName[T any]() string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	var zero T
	return tableRegistry[reflect.TypeOf(zero)]
}
