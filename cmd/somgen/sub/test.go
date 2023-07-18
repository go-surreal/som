package sub

import (
	"github.com/marcbinz/som/core/api"
	"github.com/urfave/cli/v2"
)

func Test() *cli.Command {
	return &cli.Command{
		Name:   "test",
		Action: test,
	}
}

func test(ctx *cli.Context) error {
	return api.Test(ctx.Context)
}
