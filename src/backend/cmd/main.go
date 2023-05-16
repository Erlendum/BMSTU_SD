package main

import (
	"backend/cmd/console"
	"backend/cmd/modes"
	"flag"
)

var mode = flag.String("mode", "http-server", "mode of backend")

func main() {
	flag.Parse()

	switch *mode {
	case "cli":
		app := console.App{}
		app.Run()
	case "http-server":
		app := modes.App{}
		app.Run()
	}
}
