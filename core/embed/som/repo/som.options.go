//go:build embed

package repo

import (
	"log/slog"
	"time"
)

type options struct {
	timeout   time.Duration
	logger    *slog.Logger
	readLimit int64
}

type Option func(*options)

// WithTimeout sets a custom timeout for requests.
// Note: This option is currently not used by the surrealdb.go client wrapper.
// Timeout control should be done via context.Context with timeout/deadline.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// WithLogger sets the logger.
// Note: This option is currently not used by the surrealdb.go client wrapper.
// Consider using structured logging at the application level.
func WithLogger(logger *slog.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

// WithReadLimit sets a custom read limit (in bytes) for the websocket connection.
// Note: This option is currently not used by the surrealdb.go client wrapper.
func WithReadLimit(limit int64) Option {
	return func(opts *options) {
		opts.readLimit = limit
	}
}

func applyOptions(opts []Option) *options {
	out := &options{}

	for _, opt := range opts {
		opt(out)
	}

	return out
}