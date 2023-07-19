package sdbd

import (
	"context"
	"fmt"
)

const (
	methodSignin = "signin"
	methodUse    = "use"
	methodQuery  = "query"
	methodLive   = "live"
	methodKill   = "kill"
	methodUpdate = "update"
	methodDelete = "delete"
	methodSelect = "select"
	methodCreate = "create"
	methodChange = "change"
	methodModify = "modify"
)

// signIn is a helper method for signing in a user.
func (c *Client) signIn(ctx context.Context, username, password string) error {
	res, err := c.send(ctx, Request{
		Method: methodSignin,
		Params: []any{
			map[string]string{
				"user": username,
				"pass": password,
			},
		},
	})
	if err != nil {
		return err
	}

	fmt.Println("sign_in:", res)

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

	fmt.Println("use:", res)

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

func (c *Client) Live(ctx context.Context, query string) (string, error) {
	res, err := c.send(ctx, Request{
		Method: methodQuery, // TODO: switch to methodLive once its working with it ;)
		Params: []any{
			"live " + query,
		},
	})
	if err != nil {
		return "", err
	}

	// register the returned id to receive the live updates
	// add a small queue in the background if an update might come in before the caller has registered the callback

	fmt.Println("live:", res)

	return "", nil
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

func (c *Client) Create(ctx context.Context) (interface{}, error) {
	res, err := c.send(ctx, Request{
		Method: methodCreate,
		Params: []any{
			"person",
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("create:", res)

	return "", nil
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
