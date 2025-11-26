//go:build embed

package internal

import "time"

type Timestamps struct {
	createdAt time.Time
	updatedAt time.Time
}

func NewTimestamps(createdAt time.Time, updatedAt time.Time) Timestamps {
	return Timestamps{
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (t Timestamps) CreatedAt() time.Time {
	return t.createdAt
}

func (t Timestamps) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Timestamps) SetCreatedAt(tm time.Time) {
	t.createdAt = tm
}

func (t *Timestamps) SetUpdatedAt(tm time.Time) {
	t.updatedAt = tm
}
