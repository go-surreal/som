package sub

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/gorilla"
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
	db, err := New("root", "root", "som", "default")
	if err != nil {
		return err
	}
	defer db.Close()

	// fmt.Println(db.DB.Query("relate user:k14mjkmjp0z9zulpug0a->member_of2->group:s3h3r4rmscpso3f2nv4s", nil))

	var x *string

	update, err := db.DB.Update("user:marc", map[string]any{
		"name": x,
	})
	if err != nil {
		return err
	}

	fmt.Println(update)

	live, err := db.DB.Live("select * from user")
	if err != nil {
		return err
	}

	fmt.Println("live:", live)

	var bytes []byte

	for _, val := range live.([]interface{}) {
		bytes = append(bytes, byte(val.(float64)))
	}

	uid, err := uuid.FromBytes(bytes)
	if err != nil {
		return err
	}

	fmt.Println("uid:", uid)

	query, err := db.DB.Query("select * from "+uid.String(), nil)
	if err != nil {
		return err
	}

	fmt.Println("query:", query)

	kill, err := db.DB.Kill(uid.String())
	if err != nil {
		return err
	}

	fmt.Println("kill:", kill)

	// _, err = db.Create(ctx.Context, &Data{
	// 	Key: "some key",
	// 	SomeData: SomeData{
	// 		Value: "some value",
	// 		MoreInfo: MoreInfo{
	// 			Text: "some text",
	// 		},
	// 	},
	// 	CreatedAt: time.Now(),
	// })
	// if err != nil {
	// 	return err
	// }
	//
	// rows, selectOneErr := db.Query("select * from data4 where some_data.more_info != $0", map[string]any{
	// 	"0": nil,
	// })
	// if selectOneErr != nil {
	// 	log.Fatal(selectOneErr)
	// }
	//
	// for _, row := range rows {
	// 	fmt.Println(row.Key, row.CreatedAt.Format(time.RFC3339))
	// }

	return nil
}

type Client struct {
	DB *surrealdb.DB
}

func New(username, password, namespace, database string) (*Client, error) {
	ws, err := gorilla.Create().SetTimeOut(time.Minute).Connect("ws://localhost:8020/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	db, err := surrealdb.New("<unused>", ws)
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

func (c *Client) Create(_ context.Context, data *Data) (string, error) {
	raw := toRaw(data)

	res, err := c.DB.Create("data4", raw)
	if err != nil {
		return "", err
	}

	fmt.Println(res)

	return "", nil
}

type Result struct {
	Result []*Data
	Time   string
	Status string
}

func (c *Client) Query(what string, vars map[string]any) ([]Data, error) {
	raw, err := c.DB.Query(what, vars)
	if err != nil {
		return nil, err
	}

	var res1 *Result
	err = surrealdb.Unmarshal(raw, &res1)
	if err != nil {
		return nil, err
	}

	var res2 *Data
	ok, err := surrealdb.UnmarshalRaw(raw, &res2)
	if err != nil {
		return nil, err
	}

	fmt.Println("res1:", ok, res1.Status, res1.Time, res1.Result[0].ID)
	fmt.Println("res2:", ok, res2)

	return nil, nil
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

//go:som table
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

// func fromRaw(data map[string]any) Data {
// 	fmt.Println("fromRaw:", data)
// 	return Data{
// 		ID:  data["id"].(string),
// 		Key: data["key"].(string),
// 		SomeData: SomeData{
// 			Value: data["some_data"].(map[string]any)["value"].(string),
// 			MoreInfo: MoreInfo{
// 				Text: data["some_data"].(map[string]any)["more_info"].(map[string]any)["text"].(string),
// 			},
// 		},
// 		CreatedAt: parseTime(data["created_at"].(string)),
// 	}
// }
//
// func parseTime(val string) time.Time {
// 	res, _ := time.Parse(time.RFC3339, val)
// 	return res
// }
