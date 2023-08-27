package sdbd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
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
	Error  *ResponseError  `json:"error"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LiveQueryResult struct {
	ID     []byte `json:"id"`
	Action string `json:"action"`
	Result any    `json:"result"`
}

func (c *Client) subscribe(ctx context.Context) {
	defer c.waitGroup.Done()

	ctx, cancel := context.WithCancel(ctx)

	ch := make(resultChannel[[]byte])

	go func(ch resultChannel[[]byte]) {
		defer cancel()
		defer close(ch)

		for {

			buf, err := read(ctx, c.conn)

			// typ, data, err := c.conn.Read(ctx)

			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}

				if errors.Is(err, io.EOF) || websocket.CloseStatus(err) != -1 {
					c.logger.Info("Websocket closed.")
					return
				}

				c.logger.Error("Could not read from websocket.", "error", err)
				continue
			}

			ch <- result(buf, nil)
		}
	}(ch)

	c.handleMessages(ctx, ch)
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
				c.logger.DebugContext(ctx, "Result channel closed. Stopping message handler.")
				return
			}

			data, err := result()
			if err != nil {
				c.logger.ErrorContext(ctx, "Could not get result from channel.", "error", err)
				continue
			}

			go c.handleMessage(ctx, data)
		}
	}
}

func (c *Client) handleMessage(ctx context.Context, data []byte) {
	c.waitGroup.Add(1)
	defer c.waitGroup.Done()

	var res *Response

	if err := jsonHandler.Unmarshal(data, &res); err != nil {
		// c.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal JSON")
		// return fmt.Errorf("failed to unmarshal JSON: %w", err)
		return
	}

	c.logger.InfoContext(ctx, "Received message.", "res", res)

	if res.Error != nil {
		c.logger.ErrorContext(ctx, "Received error response.", "error", res.Error)
		return
	}

	if res.ID == "" {
		c.handleLiveQuery(ctx, res)
		return
	}

	c.handleResult(ctx, res)
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

	fmt.Println("live:", string(res.Result))

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

func isPermanentError(err error) bool {
	if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
		return true
	}

	if errors.As(err, new(net.Error)) {
		return true
	}

	return errors.Is(err, io.EOF)
}
