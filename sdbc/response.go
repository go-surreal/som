package sdbc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"nhooyr.io/websocket"
	"time"
)

func (c *Client) subscribe() {
	ch := make(resultChannel[[]byte])

	c.waitGroup.Add(1)
	go func(ch resultChannel[[]byte]) {
		defer c.waitGroup.Done()

		defer close(ch)

		for {
			buf, err := c.read(c.connCtx)

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

	c.handleMessages(ch)
}

// read reads a single websocket message.
// It will reuse buffers in between calls to avoid allocations.
func (c *Client) read(ctx context.Context) (_ []byte, err error) {
	defer c.checkWebsocketConn(err)

	typ, r, err := c.conn.Reader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get reader: %w", err)
	}

	if typ != websocket.MessageText {
		return nil, fmt.Errorf("expected message of type text (%d), got %v", websocket.MessageText, typ)
	}

	b := c.buffers.Get()
	defer c.buffers.Put(b)

	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return b.Bytes(), nil
}

func (c *Client) handleMessages(resultCh resultChannel[[]byte]) {
	for {
		select {

		case <-c.connCtx.Done():
			{
				c.logger.DebugContext(c.connCtx, "Context done. Stopping message handler.")
				return
			}

		case result, more := <-resultCh:
			{
				if !more {
					c.logger.DebugContext(c.connCtx, "Result channel closed. Stopping message handler.")
					return
				}

				c.waitGroup.Add(1)
				go func() {
					defer c.waitGroup.Done()

					data, err := result()
					if err != nil {
						c.logger.ErrorContext(c.connCtx, "Could not get result from channel.", "error", err)
						return
					}

					c.handleMessage(data)
				}()
			}
		}
	}
}

func (c *Client) handleMessage(data []byte) {
	var res *response

	if err := c.jsonUnmarshal(data, &res); err != nil {
		c.logger.ErrorContext(c.connCtx, "Could not unmarshal websocket message.", "error", err)
		return
	}

	c.logger.DebugContext(c.connCtx, "Received message.", "res", res)

	if res.Error != nil {
		c.logger.ErrorContext(c.connCtx, "Received error response.", "error", res.Error)
		return
	}

	if res.ID == "" {
		c.handleLiveQuery(res)
		return
	}

	c.handleResult(res)
}

func (c *Client) handleResult(res *response) {
	outCh, ok := c.requests.get(res.ID)
	if !ok {
		c.logger.ErrorContext(c.connCtx, "Could not find pending request for ID.", "id", res.ID)
		return
	}

	select {

	case outCh <- res.Result:
		return

	case <-c.connCtx.Done():
		return

	case <-time.After(c.timeout):
		c.logger.ErrorContext(c.connCtx, "Timeout while sending result to channel.", "id", res.ID)
	}
}

func (c *Client) handleLiveQuery(res *response) {
	var rawID liveQueryID

	if err := c.jsonUnmarshal(res.Result, &rawID); err != nil {
		c.logger.ErrorContext(c.connCtx, "Could not unmarshal websocket message.", "error", err)
		return
	}

	outCh, ok := c.liveQueries.get(rawID.ID, false)
	if !ok {
		c.logger.ErrorContext(c.connCtx, "Could not find live query channel.", "id", rawID.ID)
		return
	}

	select {

	case outCh <- res.Result:
		return

	case <-c.connCtx.Done():
		return

	case <-time.After(c.timeout):
		c.logger.ErrorContext(c.connCtx, "Timeout while sending result to channel.", "id", res.ID)
	}
}
