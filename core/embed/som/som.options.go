//go:build embed

package som

import (
	"encoding/json"
	"github.com/go-surreal/sdbc"
	"log/slog"
	"time"
)

type options struct {
	jsonMarshal   JsonMarshal
	jsonUnmarshal JsonUnmarshal
	sdbc          []sdbc.Option
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

type JsonMarshal func(val any) ([]byte, error)
type JsonUnmarshal func(buf []byte, val any) error

func applyOptions(opts []Option) *options {
	out := &options{
		jsonMarshal:   json.Marshal,
		jsonUnmarshal: json.Unmarshal,
	}

	for _, opt := range opts {
		opt(out)
	}

	return out
}
