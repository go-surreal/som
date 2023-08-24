//go:build embed

package som

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/gorilla"
	"github.com/surrealdb/surrealdb.go/pkg/logger"
	"time"
)

type Database interface {
	Close()
	Create(thing string, data any) (any, error)
	Select(what string) (any, error)
	Query(statement string, vars any) (any, error)
	Update(thing string, data any) (any, error)
	Delete(what string) (any, error)
}

type Config struct {
	Address   string
	Username  string
	Password  string
	Namespace string
	Database  string
}

type ClientImpl struct {
	db Database
}

func NewClient(conf Config) (*ClientImpl, error) {
	url := conf.Address + "/rpc"

	logData, err := logger.New().Make()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %v", err)
	}

	ws, err := gorilla.Create().Logger(logData).SetTimeOut(time.Minute).Connect(url)
	if err != nil {
		return nil, fmt.Errorf("could not create websocket: %v", err)
	}

	surreal, err := surrealdb.New("<unused>", ws)
	if err != nil {
		return nil, fmt.Errorf("could not create surrealdb client: %v", err)
	}

	_, err = surreal.Signin(map[string]any{
		"user": conf.Username,
		"pass": conf.Password,
	})
	if err != nil {
		return nil, err
	}

	_, err = surreal.Use(conf.Namespace, conf.Database)
	if err != nil {
		return nil, err
	}

	rawRes, err := surreal.Query(fmt.Sprintf("DEFINE NAMESPACE %s", conf.Namespace), nil)
	if err != nil {
		return nil, err
	}

	nsRes, ok := rawRes.([]any)[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("could not create namespace: %v", rawRes)
	}

	ns, ok := nsRes["result"]
	if !ok || ns != nil {
		return nil, fmt.Errorf("could not create namespace: %v", nsRes)
	}

	rawRes, err = surreal.Query(fmt.Sprintf("DEFINE DATABASE %s", conf.Database), nil)
	if err != nil {
		return nil, err
	}

	dbRes, ok := rawRes.([]any)[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("could not create database: %v", rawRes)
	}

	db, ok := dbRes["result"]
	if !ok || db != nil {
		return nil, fmt.Errorf("could not create database: %v", dbRes)
	}

	return &ClientImpl{db: &database{DB: surreal}}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
