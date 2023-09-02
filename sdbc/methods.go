package sdbc

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"time"
)

const (
	methodSignIn = "signin"
	methodUse    = "use"
	methodQuery  = "query"
	methodLive   = "live"
	methodKill   = "kill"
	methodUpdate = "update"
	methodDelete = "delete"
	methodSelect = "select"
	methodCreate = "create"
)

const (
	nilValue = "null"
)

// signIn is a helper method for signing in a user.
func (c *Client) signIn(ctx context.Context, timeout time.Duration, username, password string) error {
	res, err := c.send(ctx,
		request{
			Method: methodSignIn,
			Params: []any{
				signInParams{
					User: username,
					Pass: password,
				},
			},
		},
		timeout,
	)
	if err != nil {
		return fmt.Errorf("could not sign in: %w", err)
	}

	c.token = string(res)

	return nil
}

// use is a method to select the namespace and table for the connection.
func (c *Client) use(ctx context.Context, timeout time.Duration, namespace, database string) error {
	res, err := c.send(ctx,
		request{
			Method: methodUse,
			Params: []any{
				namespace,
				database,
			},
		},
		timeout,
	)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if string(res) != nilValue {
		return fmt.Errorf("could not select database due to %s", string(res))
	}

	return nil
}

// Query is a convenient method for sending a query to the database.
func (c *Client) Query(ctx context.Context, timeout time.Duration, query string, vars map[string]any) ([]byte, error) {
	res, err := c.send(ctx,
		request{
			Method: methodQuery,
			Params: []any{
				query,
				vars,
			},
		},
		timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return res, nil
}

func (c *Client) Live(ctx context.Context, timeout time.Duration, query string) (<-chan []byte, error) {
	raw, err := c.send(ctx,
		request{
			Method: methodLive,
			Params: []any{
				query,
			},
		},
		timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	var res []basicResponse[string]

	if err := c.jsonUnmarshal(raw, &res); err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %w", err)
	}

	if len(res) < 1 {
		return nil, fmt.Errorf("empty response")
	}

	ch, ok := c.liveQueries.get(res[0].Result, true)
	if !ok {
		return nil, fmt.Errorf("could not get live query channel")
	}

	go func(key string) {
		<-ctx.Done()
		c.liveQueries.del(key)
	}(res[0].Result)

	return ch, nil
}

func (c *Client) Kill(ctx context.Context, timeout time.Duration, uuid string) ([]byte, error) {
	res, err := c.send(ctx,
		request{
			Method: methodKill,
			Params: []any{
				uuid,
			},
		},
		timeout,
	)
	if err != nil {
		return res, fmt.Errorf("could not send request: %w", err)
	}

	return res, nil
}

// Select a table or record from the database.
func (c *Client) Select(ctx context.Context, timeout time.Duration, thing string) ([]byte, error) {
	res, err := c.send(ctx,
		request{
			Method: methodSelect,
			Params: []any{
				thing,
			},
		},
		timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return res, nil
}

func (c *Client) Create(ctx context.Context, timeout time.Duration, thing string, data any) ([]byte, error) {
	res, err := c.send(ctx,
		request{
			Method: methodCreate,
			Params: []any{
				thing,
				data,
			},
		},
		timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return res, nil
}

// Update a table or record in the database like a PUT request.
func (c *Client) Update(ctx context.Context, timeout time.Duration, thing string, data any) ([]byte, error) {
	res, err := c.send(ctx,
		request{
			Method: methodUpdate,
			Params: []any{
				thing,
				data,
			},
		},
		timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return res, nil
}

// Delete a table or a row from the database like a DELETE request.
func (c *Client) Delete(ctx context.Context, timeout time.Duration, thing string) ([]byte, error) {
	res, err := c.send(ctx,
		request{
			Method: methodDelete,
			Params: []any{
				thing,
			},
		},
		timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return res, nil
}

//
// -- TYPES
//

type signInParams struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

//
// -- INTERNAL
//

func (c *Client) send(ctx context.Context, req request, timeout time.Duration) (_ []byte, err error) {
	defer c.checkWebsocketConn(err)

	reqID, resCh := c.requests.prepare()
	defer c.requests.cleanup(reqID)

	req.ID = reqID

	c.logger.DebugContext(ctx, "Sending request.", "request", req)

	if err := c.write(ctx, req); err != nil {
		return nil, fmt.Errorf("could not write to websocket: %w", err)
	}

	if deadline, ok := ctx.Deadline(); ok && timeout == 0 {
		timeout = time.Until(deadline)
	}

	if timeout == 0 {
		timeout = c.timeout
	}

	select {

	case <-ctx.Done():
		return nil, fmt.Errorf("context done: %w", ctx.Err())

	case <-time.After(timeout):
		return nil, fmt.Errorf("request timed out")

	case res, open := <-resCh:
		if !open {
			return nil, fmt.Errorf("channel closed")
		}

		return res, nil
	}
}

// write writes the JSON message v to c.
// It will reuse buffers in between calls to avoid allocations.
func (c *Client) write(ctx context.Context, req request) (err error) {
	defer c.checkWebsocketConn(err)

	data, err := c.jsonMarshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = c.conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	// TODO: use Writer instead of Write to stream the message?
	return nil
}
