package model

import (
	"som.test/gen/som"
)

// EventLog is a write-only ingestion sink: records are accepted to feed the
// EventSummary view and then discarded (DROP table). It exercises the
// sink→view pattern.
type EventLog struct {
	som.Sink

	Category string
	Value    float64
}

// EventSummary is a read-only view aggregating EventLog records, grouped by
// Category. Its rows are computed from the discarded EventLog writes.
type EventSummary struct {
	som.View

	Category string
	Total    int
	AvgValue float64
}
