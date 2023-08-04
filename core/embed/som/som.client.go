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

	return &ClientImpl{db: &database{DB: surreal}}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
