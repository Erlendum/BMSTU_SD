package main

import (
	"backend/cmd/modes"
)

func main() {
	app := modes.App{}
	//err := app.ParseConfig("./config", "config")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	app.Run()
}
