package sdbd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
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

func (c *Client) send(ctx context.Context, req Request) ([]byte, error) {
	reqID, resCh := c.requests.prepare()
	defer c.requests.cleanup(reqID)

	req.ID = reqID

	data, err := c.jsonMarshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %w", err)
	}

	err = c.socket.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return nil, fmt.Errorf("could not write to websocket: %w", err)
	}

	select {

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-time.After(c.timeout):
		return nil, fmt.Errorf("timeout")

	case res := <-resCh:
		return res, nil
	}
}
