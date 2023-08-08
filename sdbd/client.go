package sdbd

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

const (
	maxReadBufferSize  = 100 * 1024
	maxWriteBufferSize = 100 * 1024
)

type Client struct {
	*options

	conn      *websocket.Conn
	waitGroup sync.WaitGroup

	requests    requests
	liveQueries liveQueries
}

type Config struct {
	Address   string
	Username  string
	Password  string
	Namespace string
	Database  string
}

// NewClient creates a new client and connects to
// the database using a websocket connection.
func NewClient(ctx context.Context, conf Config, opts ...Option) (*Client, error) {
	conn, _, err := websocket.Dial(ctx, conf.Address, &websocket.DialOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not open websocket connection: %w", err)
	}

	client := &Client{
		options: applyOptions(opts),
		conn:    conn,
	}

	if client.options.readLimit > 0 {
		conn.SetReadLimit(client.options.readLimit)
	} else {
		conn.SetReadLimit(maxReadBufferSize)
	}

	client.waitGroup.Add(1)
	go client.subscribe(ctx)

	if err := client.signIn(ctx, 0, conf.Username, conf.Password); err != nil {
		return nil, fmt.Errorf("could not sign in: %v", err)
	}

	if err := client.use(ctx, 0, conf.Namespace, conf.Database); err != nil {
		return nil, fmt.Errorf("could not select database: %v", err)
	}

	return client, nil
}

func (c *Client) Close() error {
	c.logger.Info("Closing client.")

	err := c.conn.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		return fmt.Errorf("could not close websocket connection: %v", err)
	}

	defer c.requests.reset()
	defer c.liveQueries.reset()

	c.logger.Info("Waiting for goroutines to finish.")

	ch := make(chan struct{})

	go func() {
		defer close(ch)
		c.waitGroup.Wait()
	}()

	select {

	case <-ch:
		return nil

	case <-time.After(time.Second * 5):
		return fmt.Errorf("could not close websocket connection: timeout")
	}
}

// Used for RawQuery Unmarshaling
/*type RawQuery[I any] struct {
	Status string `json:"status"`
	Time   string `json:"time"`
	Result I      `json:"result"`
	Detail string `json:"detail"`
}*/
