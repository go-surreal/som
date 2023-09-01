package sdbc

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

// Get returns a buffer from the pool or
// creates a new one if the pool is empty.
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
	typ, r, err := c.Reader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get reader: %w", err)
	}

	if typ != websocket.MessageText {
		return nil, fmt.Errorf("expected message of type text (%d), got %v", websocket.MessageText, typ)
	}

	b := bpool.Get()
	defer bpool.Put(b)

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

// write writes the JSON message v to c.
// It will reuse buffers in between calls to avoid allocations.
func write(ctx context.Context, conn *websocket.Conn, req any) error {
	// defer errd.Wrap(&err, "failed to write JSON message")

	data, err := jsonHandler.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	// Using Writer instead of Write to stream the message.
	// writer, err := conn.Writer(ctx, websocket.MessageText)
	// if err != nil {
	// 	return err
	// }

	// json.Marshal cannot reuse buffers between calls as it has to return
	// a copy of the byte slice but Encoder does as it directly writes to w.
	// err = jsonHandler.NewEncoder(writer).Encode(req)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal JSON: %w", err)
	// }

	// return writer.Close()
	return nil
}
