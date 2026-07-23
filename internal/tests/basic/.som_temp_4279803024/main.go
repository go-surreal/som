//go:build som

package main

import (
	"os"
	model "som.test"
)

func main() {
	data, err := model.Definitions().ToJSON()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.Write(data)
}
