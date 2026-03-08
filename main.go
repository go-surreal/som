package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/go-surreal/som/core"
)

func main() {
	ctx := context.Background()

	app := core.Command()

	app.Authors = []any{
		"Marc Binz",
	}
	app.Copyright = "github.com/go-surreal/som"

	app.ExtraInfo = func() map[string]string {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return nil
		}

		return map[string]string{
			"GoVersion": info.GoVersion,
		}
	}

	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
