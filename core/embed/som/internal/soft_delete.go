//go:build embed

package internal

import "time"

// SoftDelete provides soft delete functionality when embedded in a model.
// The DeletedAt field is set automatically when Delete() is called.
type SoftDelete struct {
	deletedAt time.Time
}

// DeletedAt returns the time this record was soft-deleted.
func (d SoftDelete) DeletedAt() time.Time {
	return d.deletedAt
}

// IsDeleted returns true if this record has been soft-deleted.
func (d SoftDelete) IsDeleted() bool {
	return !d.deletedAt.IsZero()
}

func SetDeletedAt(d *SoftDelete, at time.Time) {
	d.deletedAt = at
}
