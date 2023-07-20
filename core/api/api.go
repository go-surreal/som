package api

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/marcbinz/som/sdbd"
	"log/slog"
	"os"
	"sync"
)

func Test(ctx context.Context) error {
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)

	slog.Info("start")

	client, err := sdbd.NewClient(ctx,
		sdbd.Config{
			Host:      "ws://localhost:8020/rpc",
			Username:  "root",
			Password:  "root",
			Namespace: "test",
			Database:  "test",
		},
		sdbd.WithLogger(slog.Default()),
		sdbd.WithJsonHandlers(sonic.Marshal, sonic.Unmarshal),
	)
	if err != nil {
		return err
	}

	live, err := client.Live(ctx, "select * from person")
	if err != nil {
		return err
	}

	go func() {
		for liveRes := range live {
			fmt.Println("live:", string(liveRes))
		}
	}()

	create, err := client.Create(ctx, "person", map[string]interface{}{
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
