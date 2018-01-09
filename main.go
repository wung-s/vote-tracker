package main

import (
	"log"

	"github.com/wung-s/gotv/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
