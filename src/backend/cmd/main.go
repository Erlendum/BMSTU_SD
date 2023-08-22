package main

import (
	"backend/cmd/cli"
	"backend/cmd/http-server"
	"flag"
)

var mode = flag.String("mode", "http-server", "mode of backend")

func main() {
	flag.Parse()

	switch *mode {
	case "cli":
		app := cli.App{}
		app.Run()
	case "http-server":
		app := http_server.App{}
		app.Run()
	}
}
