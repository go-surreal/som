package sdbd

import (
	"context"
	"encoding/json"
	"fmt"
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
	methodInsert = "insert"
	methodChange = "change"
	methodModify = "modify"
)

const (
	nilValue = "null"
)

// signIn is a helper method for signing in a user.
func (c *Client) signIn(ctx context.Context, username, password string) error {
	res, err := c.send(ctx, Request{
		Method: methodSignIn,
		Params: []any{
			signInParams{
				User: username,
				Pass: password,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not sign in: %w", err)
	}

	if string(res) != nilValue {
		return fmt.Errorf("could not sign in due to unknown response %s", string(res))
	}

	return nil
}

// use is a method to select the namespace and table for the connection.
func (c *Client) use(ctx context.Context, namespace, database string) error {
	res, err := c.send(ctx, Request{
		Method: methodUse,
		Params: []any{
			namespace,
			database,
		},
	})
	if err != nil {
		return err
	}

	if string(res) != nilValue {
		return fmt.Errorf("could not select database due to %s", string(res))
	}

	return nil
}

// Query is a convenient method for sending a query to the database.
func (c *Client) Query(ctx context.Context, query string, vars interface{}) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodQuery,
		Params: []any{
			"live " + query,
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("live:", res)

	return "", nil
}

func (c *Client) Live(ctx context.Context, query string) (<-chan []byte, error) {
	raw, err := c.send(ctx, Request{
		Method: methodQuery, // TODO: switch to methodLive once its working with it ;)
		Params: []any{
			"live " + query,
		},
	})
	if err != nil {
		return nil, err
	}

	var res []basicResponse[string]

	if err := c.jsonUnmarshal(raw, &res); err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, fmt.Errorf("no response")
	}

	ch := c.liveQueries.get(res[0].Result)

	return ch, nil
}

type basicResponse[R any] struct {
	Status string `json:"status"`
	Result R      `json:"result"`
	Time   Time   `json:"time"`
}

type Time time.Duration

func (t *Time) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	d, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	*t = Time(d)

	return nil
}

func (c *Client) Kill(ctx context.Context, uuid string) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodKill,
		Params: []any{
			uuid,
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("kill:", res)

	return "", nil
}

// Select a table or record from the database.
func (c *Client) Select(ctx context.Context) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodSelect,
		Params: []any{},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("select:", res)

	return "", nil
}

func (c *Client) Create(ctx context.Context, table string, data any) ([]byte, error) {
	res, err := c.send(ctx, Request{
		Method: methodCreate,
		Params: []any{
			table,
			data,
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Insert(ctx context.Context, table string, data []any) ([]byte, error) {
	res, err := c.send(ctx, Request{
		Method: methodInsert,
		Params: []any{
			table,
			data,
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update a table or record in the database like a PUT request.
func (c *Client) Update(ctx context.Context) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodUpdate,
		Params: []any{},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("update:", res)

	return "", nil
}

// Change a table or record in the database like a PATCH request.
func (c *Client) Change(ctx context.Context) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodChange,
		Params: []any{},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("change:", res)

	return "", nil
}

// Modify applies a series of JSONPatches to a table or record.
func (c *Client) Modify(ctx context.Context) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodModify,
		Params: []any{},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("modify:", res)

	return "", nil
}

// Delete a table or a row from the database like a DELETE request.
func (c *Client) Delete(ctx context.Context) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodDelete,
		Params: []any{},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("delete:", res)

	return "", nil
}

//
// -- HELPER
//

type signInParams struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}
