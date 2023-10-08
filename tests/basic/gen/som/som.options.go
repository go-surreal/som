// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package som

import (
	"github.com/go-surreal/sdbc"
	"log/slog"
	"time"
)

type Option func() sdbc.Option

// WithTimeout sets a custom timeout for requests.
// If not set, the default timeout is 1 minute.
func WithTimeout(timeout time.Duration) Option {
	return func() sdbc.Option {
		return sdbc.WithTimeout(timeout)
	}
}

// WithLogger sets the logger.
// If not set, no log output is created.
func WithLogger(logger *slog.Logger) Option {
	return func() sdbc.Option {
		return sdbc.WithLogger(logger)
	}
}

// WithReadLimit sets a custom read limit (in bytes) for the websocket connection.
// If not set, the default read limit is 1 MB.
func WithReadLimit(limit int64) Option {
	return func() sdbc.Option {
		return sdbc.WithReadLimit(limit)
	}
}

func applyOptions(opts []Option) []sdbc.Option {
	var out []sdbc.Option

	for _, opt := range opts {
		out = append(out, opt())
	}

	return out
}
