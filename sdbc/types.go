package sdbc

import (
	"encoding/json"
	"fmt"
	"time"
)

type request struct {
	ID     string        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type response struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *responseError  `json:"error"`
}

type responseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type liveQueryID struct {
	ID string `json:"id"`
}

//
// -- INTERNAL
//

type basicResponse[R any] struct {
	Status string   `json:"status"`
	Result R        `json:"result"`
	Time   duration `json:"time"`
}

type duration time.Duration

func (t *duration) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("could not unmarshal duration: %w", err)
	}

	d, err := time.ParseDuration(str)
	if err != nil {
		return fmt.Errorf("could not parse duration: %w", err)
	}

	*t = duration(d)

	return nil
}

func result[T any](t T, err error) resultFunc[T] {
	return func() (T, error) {
		return t, err
	}
}

type resultFunc[T any] func() (T, error)

type resultChannel[T any] chan resultFunc[T]
