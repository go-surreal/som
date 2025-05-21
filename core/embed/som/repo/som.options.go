//go:build embed

package repo

import (
	"github.com/go-surreal/sdbc"
	"log/slog"
	"time"
)

type options struct {
	sdbc []sdbc.Option
}

type Option func(*options)

// WithTimeout sets a custom timeout for requests.
// If not set, the default timeout is 1 minute.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.sdbc = append(opts.sdbc, sdbc.WithTimeout(timeout))
	}
}

// WithLogger sets the logger.
// If not set, no log output is created.
func WithLogger(logger *slog.Logger) Option {
	return func(opts *options) {
		opts.sdbc = append(opts.sdbc, sdbc.WithLogger(logger))
	}
}

// WithReadLimit sets a custom read limit (in bytes) for the websocket connection.
// If not set, the default read limit is 1 MB.
func WithReadLimit(limit int64) Option {
	return func(opts *options) {
		opts.sdbc = append(opts.sdbc, sdbc.WithReadLimit(limit))
	}
}

func applyOptions(opts []Option) *options {
	out := &options{}

	for _, opt := range opts {
		opt(out)
	}

	return out
}
