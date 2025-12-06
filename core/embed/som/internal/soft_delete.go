//go:build embed

package internal

import "time"

// SoftDelete provides soft delete functionality when embedded in a model.
// The DeletedAt field is set automatically when Delete() is called.
type SoftDelete struct {
	DeletedAt *time.Time
}

// IsDeleted returns true if this record has been soft-deleted.
func (s SoftDelete) IsDeleted() bool {
	return s.DeletedAt != nil && !s.DeletedAt.IsZero()
}
