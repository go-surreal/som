//go:build embed

package som

import (
	"context"
	"errors"
	"time"
)

// ErrCacheSizeLimitExceeded is returned when eager cache would exceed the max size.
var ErrCacheSizeLimitExceeded = errors.New("cache size limit exceeded")

// CacheMode specifies the caching strategy.
type CacheMode int

const (
	// CacheModeLazy loads records on-demand as they are requested.
	CacheModeLazy CacheMode = iota
	// CacheModeEager loads all records on first Read() call.
	CacheModeEager
)

// CacheOptions configures cache behavior.
type CacheOptions struct {
	Mode    CacheMode
	TTL     time.Duration
	MaxSize int
}

// DefaultMaxSize is the default maximum number of records for eager cache.
const DefaultMaxSize = 1000

// CacheOption is a functional option for configuring cache behavior.
type CacheOption func(*CacheOptions)

// Lazy sets the cache to lazy loading mode (default).
// Records are fetched from the database on first access and cached.
func Lazy() CacheOption {
	return func(o *CacheOptions) {
		o.Mode = CacheModeLazy
	}
}

// Eager sets the cache to eager loading mode.
// All records are loaded on first Read() call.
func Eager() CacheOption {
	return func(o *CacheOptions) {
		o.Mode = CacheModeEager
	}
}

// WithTTL sets the time-to-live for cache entries.
// Entries expire after the given duration.
func WithTTL(d time.Duration) CacheOption {
	return func(o *CacheOptions) {
		o.TTL = d
	}
}

// WithMaxSize sets the maximum number of records for eager cache.
// Default is 1000. If the table has more records than maxSize,
// ErrCacheSizeLimitExceeded is returned.
func WithMaxSize(n int) CacheOption {
	return func(o *CacheOptions) {
		o.MaxSize = n
	}
}

// cacheKey is a unique context key for each table's cache.
type cacheKey struct {
	table string
}

// cacheOptionsKey is a unique context key for cache options.
type cacheOptionsKey struct {
	table string
}

// WithCache enables caching for the specified model type.
// The cache is stored in the context and used by subsequent Read calls.
//
// Example:
//
//	ctx = som.WithCache[model.Group](ctx)                              // lazy (default)
//	ctx = som.WithCache[model.Group](ctx, som.Lazy())                  // explicit lazy
//	ctx = som.WithCache[model.Group](ctx, som.Eager())                 // eager, load on first read
//	ctx = som.WithCache[model.Group](ctx, som.WithTTL(5*time.Minute))  // with expiration
//	ctx = som.WithCache[model.Group](ctx, som.Eager(), som.WithMaxSize(5000))
func WithCache[T any](ctx context.Context, opts ...CacheOption) context.Context {
	tableName := getTableName[T]()
	if tableName == "" {
		return ctx
	}

	options := &CacheOptions{
		Mode:    CacheModeLazy,
		MaxSize: DefaultMaxSize,
	}
	for _, opt := range opts {
		opt(options)
	}

	key := cacheKey{table: tableName}
	optsKey := cacheOptionsKey{table: tableName}

	ctx = context.WithValue(ctx, key, struct{}{})
	ctx = context.WithValue(ctx, optsKey, options)

	return ctx
}

// DropCache removes the cache for the specified model type from context.
// Subsequent Read calls using the returned context will query the database directly.
//
// Example:
//
//	ctx = som.DropCache[model.Group](ctx)
func DropCache[T any](ctx context.Context) context.Context {
	tableName := getTableName[T]()
	if tableName == "" {
		return ctx
	}

	key := cacheKey{table: tableName}
	optsKey := cacheOptionsKey{table: tableName}

	ctx = context.WithValue(ctx, key, nil)
	ctx = context.WithValue(ctx, optsKey, nil)

	return ctx
}

// CacheEnabled returns true if caching is enabled for the specified table.
func CacheEnabled(ctx context.Context, table string) bool {
	key := cacheKey{table: table}
	v := ctx.Value(key)
	return v != nil
}

// GetCacheOptions retrieves cache options for the specified table from context.
func GetCacheOptions(ctx context.Context, table string) *CacheOptions {
	optsKey := cacheOptionsKey{table: table}
	if opts, ok := ctx.Value(optsKey).(*CacheOptions); ok {
		return opts
	}
	return nil
}
