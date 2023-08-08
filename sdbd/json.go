package sdbd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"nhooyr.io/websocket"
	"sync"
)

var jsonHandler = sonic.ConfigFastest

type myPool struct {
	sync.Pool
}

var bpool myPool

// Get returns a buffer from the pool or creates a new one if
// the pool is empty.
func (p *myPool) Get() *bytes.Buffer {
	b := p.Pool.Get()
	if b == nil {
		return &bytes.Buffer{}
	}
	return b.(*bytes.Buffer)
}

// Put returns a buffer into the pool.
func (p *myPool) Put(b *bytes.Buffer) {
	b.Reset()
	p.Pool.Put(b)
}

// read reads a JSON message from c into v.
// It will reuse buffers in between calls to avoid allocations.
func read(ctx context.Context, c *websocket.Conn) ([]byte, error) {
	// defer errd.Wrap(&err, "failed to read JSON message")

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageText {
		// c.logger.Error("Received message of unsupported type, expected text. Skipping.")
		return nil, fmt.Errorf("received message of unsupported type, expected text")
	}

	b := bpool.Get()
	defer bpool.Put(b)

	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	// err = jsonHandler.Unmarshal(b.Bytes(), v)
	// if err != nil {
	// 	c.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal JSON")
	// 	return fmt.Errorf("failed to unmarshal JSON: %w", err)
	// }

	return b.Bytes(), nil
}

// write writes the JSON message v to c.
// It will reuse buffers in between calls to avoid allocations.
func write(ctx context.Context, c *websocket.Conn, v interface{}) error {
	// defer errd.Wrap(&err, "failed to write JSON message")

	// Using Writer instead of Write to stream the message.
	w, err := c.Writer(ctx, websocket.MessageText)
	if err != nil {
		return err
	}

	// json.Marshal cannot reuse buffers between calls as it has to return
	// a copy of the byte slice but Encoder does as it directly writes to w.
	err = jsonHandler.NewEncoder(w).Encode(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return w.Close()
}
