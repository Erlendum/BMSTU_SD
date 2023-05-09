package main

import (
	"backend/cmd/modes"
)

func main() {
	app := modes.App{}
	//app := console.App{}

	app.Run()
}
