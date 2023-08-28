package api

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/marcbinz/som/sdbc"
	"log/slog"
	"os"
	"sync"
)

func Test(ctx context.Context) error {
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	slog.Info("start")

	client, err := sdbc.NewClient(ctx,
		sdbc.Config{
			Address:   "ws://localhost:8020/rpc",
			Username:  "root",
			Password:  "root",
			Namespace: "test",
			Database:  "test",
		},
		sdbc.WithLogger(slog.Default()),
		sdbc.WithJsonHandlers(sonic.Marshal, sonic.Unmarshal),
	)
	if err != nil {
		return err
	}

	live, err := client.Live(ctx, 0, "select * from person")
	if err != nil {
		return err
	}

	go func() {
		for liveRes := range live {
			fmt.Println("live:", string(liveRes))
		}
	}()

	create, err := client.Create(ctx, 0, "person", map[string]interface{}{
		"name": "some",
	})
	if err != nil {
		return err
	}

	fmt.Println("create:", string(create))

	/*liveQuery := res.([]any)[0].(map[string]any)["result"].(string)

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

	fmt.Println("res:", res)*/

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	return nil
}

//
// -- X
//
