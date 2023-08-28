package sdbc

import (
	"encoding/json"
	"log/slog"
	"time"
)

type options struct {
	timeout       time.Duration
	logger        *slog.Logger
	jsonMarshal   JsonMarshal
	jsonUnmarshal JsonUnmarshal
	readLimit     int64
}

type Option func(*options)

func WithTimeout(timeout time.Duration) Option {
	return func(c *options) {
		c.timeout = timeout
	}
}

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
	out := &options{
		timeout: time.Minute,
	}

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
