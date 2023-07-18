package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"sync"
)

type Request struct {
	Id     string        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type Response struct {
	Id     string `json:"id"`
	Action string `json:"action"`
	Result any    `json:"result"`
}

func Test(ctx context.Context) error {

	conn, err := newConnection(ctx, "ws://localhost:8020/rpc")
	if err != nil {
		return err
	}

	req := Request{
		Method: "signin",
		Params: []interface{}{
			map[string]interface{}{
				"user": "root",
				"pass": "root",
			},
		},
	}

	res, err := conn.Send(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println("res:", res)

	req = Request{
		Method: "use",
		Params: []interface{}{
			"test",
			"test",
		},
	}

	res, err = conn.Send(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println("res:", res)

	req = Request{
		Method: "live",
		Params: []interface{}{
			"select * from person",
		},
	}

	res, err = conn.Send(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println("res:", res)

	liveQuery := res.(string)

	req = Request{
		Method: "create",
		Params: []interface{}{
			"person",
			map[string]interface{}{
				"name": "some",
			},
		},
	}

	res, err = conn.Send(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println("res:", res)

	req = Request{
		Method: "kill",
		Params: []interface{}{
			liveQuery,
		},
	}

	res, err = conn.Send(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println("res:", res)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	return nil
}

//
// -- X
//

type Connection struct {
	ws       *websocket.Conn
	requests map[string]chan any
	mu       sync.Mutex
}

func newConnection(ctx context.Context, url string) (*Connection, error) {
	ws, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		return nil, fmt.Errorf("could not open websocket connection: %v", err)
	}

	requests := map[string]chan any{}

	go x(ctx, ws, requests)

	return &Connection{
		ws:       ws,
		requests: requests,
	}, nil
}

func x(ctx context.Context, ws *websocket.Conn, requests map[string]chan any) {
	for {
		typ, data, err := ws.Read(ctx)
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		if typ != websocket.MessageText {
			fmt.Println("read: not text")
			return
		}

		var res Response

		if err := json.Unmarshal(data, &res); err != nil {
			fmt.Println("read:", string(data))
			fmt.Println("unmarshal:", err)
			return
		}

		fmt.Println("read:", string(data))

		ch, ok := requests[res.Id]
		if !ok {
			fmt.Println("read: no request")
			return
		}

		ch <- res.Result
	}
}

func (c *Connection) Close() error {
	return c.ws.Close(websocket.StatusNormalClosure, "done")
}

func (c *Connection) prepareRequest() (uuid.UUID, chan any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := uuid.New()
	ch := make(chan any)

	c.requests[id.String()] = ch

	return id, ch
}

func (c *Connection) removeRequest(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.requests[id.String()]; ok {
		close(ch)
	}

	delete(c.requests, id.String())
}

func (c *Connection) Send(ctx context.Context, req Request) (any, error) {
	reqID, resChan := c.prepareRequest()
	defer c.removeRequest(reqID)

	req.Id = reqID.String()

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = c.ws.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return nil, err
	}

	return <-resChan, nil
}
