package sub

import (
	"context"
	"fmt"
	"github.com/surrealdb/surrealdb.go"
	"github.com/urfave/cli/v2"
	"time"
)

func Surreal() *cli.Command {
	return &cli.Command{
		Name:   "surreal",
		Action: surreal,
	}
}

func surreal(ctx *cli.Context) error {
	db, err := New("root", "root", "sdb", "default")
	if err != nil {
		return err
	}
	defer db.Close()

	id, err := db.Create(ctx.Context, &Data{
		Key: "some key",
		SomeData: SomeData{
			Value: "some value",
			MoreInfo: MoreInfo{
				Text: "some text",
			},
		},
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Println("id:", id)

	rows, selectOneErr := db.Query("select * from data4 where some_data.more_info != $0", map[string]any{
		"0": nil,
	})
	if selectOneErr != nil {
		panic(selectOneErr)
	}

	for _, row := range rows {
		fmt.Println(row.Key, row.CreatedAt.Format(time.RFC3339))
	}

	return nil
}

type Client struct {
	DB *surrealdb.DB
}

func New(username, password, namespace, database string) (*Client, error) {
	db, err := surrealdb.New("ws://localhost:8010/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	_, err = db.Signin(map[string]any{
		"user": username,
		"pass": password,
	})
	if err != nil {
		return nil, err
	}

	_, err = db.Use(namespace, database)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}

func (c *Client) Close() {
	c.DB.Close()
}

func (c *Client) Create(ctx context.Context, data *Data) (string, error) {
	raw := toRaw(data)

	res, err := c.DB.Create("data4", raw)
	if err != nil {
		return "", err
	}

	fmt.Println(res)

	return "", nil
}

func (c *Client) Query(what string, vars map[string]any) ([]Data, error) {
	res, err := c.DB.Query(what, vars)
	if err != nil {
		return nil, err
	}

	castedRes := res.([]any)

	castedFirst := castedRes[0].(map[string]any)

	fmt.Println("query output:", castedFirst["status"], castedFirst["time"])

	castedRows := castedFirst["result"].([]any)

	var out []Data
	for _, row := range castedRows {
		castedRow := row.(map[string]any)
		res = append(out, fromRaw(castedRow))
	}

	return out, nil
}

type Data struct {
	ID        string
	Key       string
	SomeData  SomeData
	CreatedAt time.Time
}

type SomeData struct {
	Value    string
	MoreInfo MoreInfo
}

//go:sdb table
type MoreInfo struct {
	Text string
}

func toRaw(data *Data) map[string]any {
	return map[string]any{
		"key": data.Key,
		"some_data": map[string]any{
			"value": data.SomeData.Value,
			"more_info": map[string]any{
				"text": data.SomeData.MoreInfo.Text,
			},
		},
		"created_at": data.CreatedAt, // .Format(time.RFC3339)
	}
}

func fromRaw(data map[string]any) Data {
	fmt.Println("fromRaw:", data)
	return Data{
		ID:  data["id"].(string),
		Key: data["key"].(string),
		SomeData: SomeData{
			Value: data["some_data"].(map[string]any)["value"].(string),
			MoreInfo: MoreInfo{
				Text: data["some_data"].(map[string]any)["more_info"].(map[string]any)["text"].(string),
			},
		},
		CreatedAt: parseTime(data["created_at"].(string)),
	}
}

func parseTime(val string) time.Time {
	res, _ := time.Parse(time.RFC3339, val)
	return res
}
