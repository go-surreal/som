//go:build embed

package lib

import "errors"

// PaginateConfig holds the configuration for cursor-based pagination.
type PaginateConfig struct {
	// First requests the first N items (forward pagination).
	First int

	// After is the cursor to start after (forward pagination).
	After string

	// Last requests the last N items (backward pagination).
	Last int

	// Before is the cursor to end before (backward pagination).
	Before string

	// IncludeTotalCount enables fetching the total count of matching records.
	IncludeTotalCount bool

	// AccuratePageInfo enables extra queries to accurately determine
	// HasPreviousPage and HasNextPage. Without this, HasPreviousPage is
	// assumed false on the first page (when no cursor is provided).
	AccuratePageInfo bool
}

// PaginateOption is a functional option for configuring pagination.
type PaginateOption func(*PaginateConfig)

// First sets forward pagination to return the first n items.
func First(n int) PaginateOption {
	return func(c *PaginateConfig) {
		c.First = n
	}
}

// After sets the cursor to start after for forward pagination.
func After(cursor string) PaginateOption {
	return func(c *PaginateConfig) {
		c.After = cursor
	}
}

// Last sets backward pagination to return the last n items.
func Last(n int) PaginateOption {
	return func(c *PaginateConfig) {
		c.Last = n
	}
}

// Before sets the cursor to end before for backward pagination.
func Before(cursor string) PaginateOption {
	return func(c *PaginateConfig) {
		c.Before = cursor
	}
}

// WithTotalCount enables fetching the total count of matching records.
// This requires an additional COUNT query.
func WithTotalCount() PaginateOption {
	return func(c *PaginateConfig) {
		c.IncludeTotalCount = true
	}
}

// WithAccuratePageInfo enables extra queries to accurately determine
// HasPreviousPage and HasNextPage. Without this option, HasPreviousPage
// is assumed false on the first page (when no cursor is provided), which
// may be inaccurate if there are records before the current page.
// This adds 1-2 extra queries but provides accurate pagination info.
func WithAccuratePageInfo() PaginateOption {
	return func(c *PaginateConfig) {
		c.AccuratePageInfo = true
	}
}

// IsBackward returns true if this is backward pagination.
func (c *PaginateConfig) IsBackward() bool {
	return c.Last > 0 || c.Before != ""
}

// Limit returns the page size (First or Last).
func (c *PaginateConfig) Limit() int {
	if c.First > 0 {
		return c.First
	}
	return c.Last
}

// Cursor returns the cursor (After or Before).
func (c *PaginateConfig) Cursor() string {
	if c.After != "" {
		return c.After
	}
	return c.Before
}

// Validate checks that the pagination configuration is valid.
func (c *PaginateConfig) Validate() error {
	if c.First > 0 && c.Last > 0 {
		return errors.New("cannot use both First and Last")
	}
	if c.After != "" && c.Before != "" {
		return errors.New("cannot use both After and Before")
	}
	if c.First == 0 && c.Last == 0 {
		return errors.New("must specify First or Last")
	}
	if c.First < 0 || c.Last < 0 {
		return errors.New("page size must be positive")
	}
	// Forward: First with optional After
	// Backward: Last with optional Before
	if c.First > 0 && c.Before != "" {
		return errors.New("cannot use Before with First (use After for forward pagination)")
	}
	if c.Last > 0 && c.After != "" {
		return errors.New("cannot use After with Last (use Before for backward pagination)")
	}
	return nil
}
