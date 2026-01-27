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

func SetCreatedAt(t *Timestamps, at time.Time) {
	t.createdAt = at
}

func SetUpdatedAt(t *Timestamps, at time.Time) {
	t.updatedAt = at
}
