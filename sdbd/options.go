package sdbd

import (
	"encoding/json"
	"log/slog"
)

type options struct {
	logger        *slog.Logger
	jsonMarshal   JsonMarshal
	jsonUnmarshal JsonUnmarshal
	readLimit     int64
}

type Option func(*options)

func WithLogger(logger *slog.Logger) Option {
	return func(c *options) {
		c.logger = logger
	}
}

func WithJsonHandlers(marshal JsonMarshal, unmarshal JsonUnmarshal) Option {
	return func(c *options) {
		c.jsonMarshal = marshal
		c.jsonUnmarshal = unmarshal
	}
}

func WithReadLimit(limit int64) Option {
	return func(c *options) {
		c.readLimit = limit
	}
}

type JsonMarshal func(val any) ([]byte, error)
type JsonUnmarshal func(buf []byte, val any) error

func applyOptions(opts []Option) *options {
	out := &options{}

	for _, opt := range opts {
		opt(out)
	}

	if out.logger == nil {
		out.logger = slog.Default() // TODO: empty impl instead?
	}

	if out.jsonMarshal == nil {
		out.jsonMarshal = json.Marshal
	}

	if out.jsonUnmarshal == nil {
		out.jsonUnmarshal = json.Unmarshal
	}

	return out
}
