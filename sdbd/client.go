package sdbd

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
)

type Client struct {
	*options

	socket      *websocket.Conn
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

func NewClient(ctx context.Context, conf Config, opts ...Option) (*Client, error) {
	ws, _, err := websocket.Dial(ctx, conf.Address, &websocket.DialOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		return nil, fmt.Errorf("could not open websocket connection: %w", err)
	}

	client := &Client{
		options: applyOptions(opts),
		socket:  ws,
	}

	if client.options.readLimit > 0 {
		ws.SetReadLimit(client.options.readLimit)
	}

	client.subscribe(ctx)

	if err := client.signIn(ctx, conf.Username, conf.Password); err != nil {
		return nil, fmt.Errorf("could not sign in: %v", err)
	}

	if err := client.use(ctx, conf.Namespace, conf.Database); err != nil {
		return nil, fmt.Errorf("could not select database: %v", err)
	}

	return client, nil
}

func (c *Client) Close() error {
	err := c.socket.Close(websocket.StatusNormalClosure, "done")
	if err != nil {
		return fmt.Errorf("could not close websocket connection: %v", err)
	}

	return nil
}

// Used for RawQuery Unmarshaling
/*type RawQuery[I any] struct {
	Status string `json:"status"`
	Time   string `json:"time"`
	Result I      `json:"result"`
	Detail string `json:"detail"`
}*/
