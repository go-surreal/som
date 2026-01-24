//go:build embed

package repo

import (
	"sync"
	"time"
)

// cache holds cached records for a specific model type.
// N is the node/model type.
type cache[N any] struct {
	mu      sync.RWMutex
	data    map[string]*cacheEntry[N]
	mode    cacheMode
	loaded  bool          // true if all records have been loaded (eager mode)
	maxSize int           // max records for eager mode
	ttl     time.Duration // time-to-live for entries (0 = no expiration)
}

// cacheEntry wraps a cached record with optional expiration time.
type cacheEntry[N any] struct {
	record    *N
	expiresAt time.Time // zero value means no expiration
}

// cacheMode specifies the caching strategy.
type cacheMode int

const (
	cacheModeLazy cacheMode = iota
	cacheModeEager
)

// newCache creates a new empty cache with the given options.
func newCache[N any](mode cacheMode, ttl time.Duration, maxSize int) *cache[N] {
	return &cache[N]{
		data:    make(map[string]*cacheEntry[N]),
		mode:    mode,
		loaded:  false,
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// newCacheWithAll creates a new cache pre-populated with all records (eager mode).
func newCacheWithAll[N any](records []*N, idFunc func(*N) string, ttl time.Duration) *cache[N] {
	c := &cache[N]{
		data:   make(map[string]*cacheEntry[N], len(records)),
		mode:   cacheModeEager,
		loaded: true,
		ttl:    ttl,
	}

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	for _, record := range records {
		if record == nil {
			continue
		}
		id := idFunc(record)
		if id == "" {
			continue
		}
		c.data[id] = &cacheEntry[N]{
			record:    record,
			expiresAt: expiresAt,
		}
	}
	return c
}

// get retrieves a record from the cache.
// Returns (record, true) if found and not expired, (nil, false) otherwise.
func (c *cache[N]) get(id string) (*N, bool) {
	c.mu.RLock()
	entry, ok := c.data[id]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	// No expiration configured
	if entry.expiresAt.IsZero() {
		return entry.record, true
	}

	// Not yet expired
	if time.Now().Before(entry.expiresAt) {
		return entry.record, true
	}

	// Possibly expired - acquire write lock and re-check
	c.mu.Lock()
	defer c.mu.Unlock()

	currentEntry, ok := c.data[id]
	if !ok {
		return nil, false
	}

	// Re-check expiration on current entry (may have been updated)
	if !currentEntry.expiresAt.IsZero() && time.Now().After(currentEntry.expiresAt) {
		delete(c.data, id)
		return nil, false
	}

	return currentEntry.record, true
}

// set stores a record in the cache.
func (c *cache[N]) set(id string, record *N) {
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

// isEager returns true if the cache was configured for eager loading.
func (c *cache[N]) isEager() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mode == cacheModeEager
}

// isLoaded returns true if all records have been loaded (for eager mode).
func (c *cache[N]) isLoaded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.loaded
}

// setLoaded marks the cache as fully loaded.
func (c *cache[N]) setLoaded() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.loaded = true
}

// maxSize returns the maximum size configured for this cache.
func (c *cache[N]) getMaxSize() int {
	return c.maxSize
}

// mode returns the cache mode.
func (c *cache[N]) getMode() cacheMode {
	return c.mode
}
