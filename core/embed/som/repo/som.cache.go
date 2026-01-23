//go:build embed

package repo

import (
	"sync"
	"time"
)

// Cache holds cached records for a specific model type.
// N is the node/model type.
type Cache[N any] struct {
	mu      sync.RWMutex
	data    map[string]*cacheEntry[N]
	mode    CacheMode
	loaded  bool          // true if all records have been loaded (eager mode)
	maxSize int           // max records for eager mode
	ttl     time.Duration // time-to-live for entries (0 = no expiration)
}

// cacheEntry wraps a cached record with optional expiration time.
type cacheEntry[N any] struct {
	record    *N
	expiresAt time.Time // zero value means no expiration
}

// CacheMode specifies the caching strategy.
type CacheMode int

const (
	CacheModeLazy CacheMode = iota
	CacheModeEager
)

// newCache creates a new empty cache with the given options.
func newCache[N any](mode CacheMode, ttl time.Duration, maxSize int) *Cache[N] {
	return &Cache[N]{
		data:    make(map[string]*cacheEntry[N]),
		mode:    mode,
		loaded:  false,
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// newCacheWithAll creates a new cache pre-populated with all records (eager mode).
func newCacheWithAll[N any](records []*N, idFunc func(*N) string, ttl time.Duration) *Cache[N] {
	c := &Cache[N]{
		data:   make(map[string]*cacheEntry[N], len(records)),
		mode:   CacheModeEager,
		loaded: true,
		ttl:    ttl,
	}

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	for _, record := range records {
		if record != nil {
			c.data[idFunc(record)] = &cacheEntry[N]{
				record:    record,
				expiresAt: expiresAt,
			}
		}
	}
	return c
}

// Get retrieves a record from the cache.
// Returns (record, true) if found and not expired, (nil, false) otherwise.
func (c *Cache[N]) Get(id string) (*N, bool) {
	c.mu.RLock()
	entry, ok := c.data[id]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.data, id)
		c.mu.Unlock()
		return nil, false
	}

	return entry.record, true
}

// Set stores a record in the cache.
func (c *Cache[N]) Set(id string, record *N) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if c.ttl > 0 {
		expiresAt = time.Now().Add(c.ttl)
	}

	c.data[id] = &cacheEntry[N]{
		record:    record,
		expiresAt: expiresAt,
	}
}

// IsEager returns true if the cache was configured for eager loading.
func (c *Cache[N]) IsEager() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mode == CacheModeEager
}

// IsLoaded returns true if all records have been loaded (for eager mode).
func (c *Cache[N]) IsLoaded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.loaded
}

// SetLoaded marks the cache as fully loaded.
func (c *Cache[N]) SetLoaded() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.loaded = true
}

// MaxSize returns the maximum size configured for this cache.
func (c *Cache[N]) MaxSize() int {
	return c.maxSize
}

// Mode returns the cache mode.
func (c *Cache[N]) Mode() CacheMode {
	return c.mode
}
