package sdbd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type requests struct {
	store sync.Map
}

func (r *requests) prepare() (string, <-chan []byte) {
	key := uuid.New()
	ch := make(chan []byte)

	r.store.Store(key.String(), ch)

	return key.String(), ch
}

func (r *requests) get(key string) (chan<- []byte, bool) {
	val, ok := r.store.Load(key)
	if !ok {
		return nil, false
	}

	return val.(chan []byte), true
}

func (r *requests) cleanup(key string) {
	if ch, ok := r.store.LoadAndDelete(key); ok {
		close(ch.(chan []byte))
	}
}

func (r *requests) reset() {
	r.store.Range(func(key, ch any) bool {
		close(ch.(chan []byte))
		r.store.Delete(key)
		return true
	})
}

func (c *Client) send(ctx context.Context, req Request, timeout time.Duration) ([]byte, error) {
	reqID, resCh := c.requests.prepare()
	defer c.requests.cleanup(reqID)

	req.ID = reqID

	c.logger.InfoContext(ctx, "Sending request.", "request", req)

	if err := write(ctx, c.conn, req); err != nil {
		return nil, fmt.Errorf("could not write to websocket: %w", err)
	}

	if deadline, ok := ctx.Deadline(); ok && timeout == 0 {
		timeout = time.Until(deadline) // TODO: okay?
	}

	if timeout == 0 {
		timeout = c.timeout
	}

	select {

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-time.After(timeout):
		return nil, fmt.Errorf("request timed out")

	case res, open := <-resCh:
		if !open {
			return nil, fmt.Errorf("channel closed")
		}

		return res, nil
	}
}
