package sdbd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"time"
)

type Request struct {
	ID     string        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type Response struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result"`
}

type LiveQueryResult struct {
	ID     []byte `json:"id"`
	Action string `json:"action"`
	Result any    `json:"result"`
}

func (c *Client) subscribe(ctx context.Context) {
	ch := make(resultChannel[[]byte])

	go func(ch resultChannel[[]byte]) {
		defer close(ch)

		for {
			typ, data, err := c.socket.Read(ctx)

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}

			if err != nil {
				c.logger.Error("Could not read from websocket.", "error", err)
				continue
			}

			if typ != websocket.MessageText {
				c.logger.Error("Received message of unsupported type, expected text. Skipping.")
				continue
			}

			ch <- result(data, nil)
		}
	}(ch)

	go c.handleMessages(ctx, ch)
}

func result[T any](t T, err error) resultFunc[T] {
	return func() (T, error) {
		return t, err
	}
}

type resultFunc[T any] func() (T, error)

type resultChannel[T any] chan resultFunc[T]

func (c *Client) handleMessages(ctx context.Context, resultCh resultChannel[[]byte]) {
	for {
		select {

		case <-ctx.Done():
			c.logger.DebugContext(ctx, "Context done. Stopping message handler.")
			return

		case result, more := <-resultCh:

			if !more {
				fmt.Println("no more")
				return
			}

			data, err := result()
			if err != nil {
				c.logger.ErrorContext(ctx, "Could not get result from channel.", "error", err)
				continue
			}

			var res *Response

			if err := c.jsonUnmarshal(data, &res); err != nil {
				c.logger.ErrorContext(ctx, "Could not unmarshal websocket message.", "error", err)
				continue
			}

			if res.ID == "" {
				c.handleLiveQuery(ctx, res)
				continue
			}

			// Do not block here, as the channel might be full.
			go c.handleResult(ctx, res)
		}
	}
}

func (c *Client) handleResult(ctx context.Context, res *Response) {
	outCh, ok := c.requests.get(res.ID)
	if !ok {
		c.logger.ErrorContext(ctx, "Could not find pending request for ID.", "id", res.ID)
		return
	}

	select {

	case outCh <- res.Result:
		return

	case <-time.After(c.timeout):
		c.logger.ErrorContext(ctx, "Timeout while sending result to channel.", "id", res.ID)
	}
}

func (c *Client) handleLiveQuery(ctx context.Context, res *Response) {
	var result LiveQueryResult

	if err := c.jsonUnmarshal(res.Result, &result); err != nil {
		c.logger.ErrorContext(ctx, "Could not unmarshal websocket message.", "error", err)
		return
	}

	uid, _ := uuid.FromBytes(result.ID) // TODO: will only work while serialization issue exists

	outCh := c.liveQueries.get(uid.String())

	select {

	case outCh <- res.Result:
		return

	case <-time.After(c.timeout):
		c.logger.ErrorContext(ctx, "Timeout while sending result to channel.", "id", res.ID)
	}
}
