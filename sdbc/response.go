package sdbc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"nhooyr.io/websocket"
	"time"
)

func (c *Client) subscribe(ctx context.Context) {
	c.waitGroup.Add(1)
	defer c.waitGroup.Done()

	ctx, cancel := context.WithCancel(ctx)

	ch := make(resultChannel[[]byte])

	go func(ch resultChannel[[]byte]) {
		defer cancel()
		defer close(ch)

		for {
			buf, err := c.read(ctx)

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

// read reads a JSON message from c into v.
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

	// err = jsonHandler.Unmarshal(b.Bytes(), v)
	// if err != nil {
	// 	c.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal JSON")
	// 	return fmt.Errorf("failed to unmarshal JSON: %w", err)
	// }

	return b.Bytes(), nil
}

func (c *Client) handleMessages(ctx context.Context, resultCh resultChannel[[]byte]) {
	for {
		select {

		case <-ctx.Done():
			{
				c.logger.DebugContext(ctx, "Context done. Stopping message handler.")
				return
			}

		case result, more := <-resultCh:
			{
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
}

func (c *Client) handleMessage(ctx context.Context, data []byte) {
	c.waitGroup.Add(1)
	defer c.waitGroup.Done()

	var res *response

	if err := c.jsonUnmarshal(data, &res); err != nil {
		return
	}

	fmt.Println("res:", fmt.Sprintf("%+v", string(data)))

	c.logger.DebugContext(ctx, "Received message.", "res", res)

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

func (c *Client) handleResult(ctx context.Context, res *response) {
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

func (c *Client) handleLiveQuery(ctx context.Context, res *response) {
	var rawID liveQueryID

	if err := c.jsonUnmarshal(res.Result, &rawID); err != nil {
		c.logger.ErrorContext(ctx, "Could not unmarshal websocket message.", "error", err)
		return
	}

	outCh, ok := c.liveQueries.get(rawID.ID, false)
	if !ok {
		c.logger.ErrorContext(ctx, "Could not find live query channel.", "id", rawID.ID)
		return
	}

	select {

	case outCh <- res.Result:
		return

	case <-time.After(c.timeout):
		c.logger.ErrorContext(ctx, "Timeout while sending result to channel.", "id", res.ID)
	}
}
