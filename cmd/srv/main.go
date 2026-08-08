package main

import (
	"flag"
	"fmt"
	"os"

	"srv.exe.dev/srv"
)

var flagListenAddr = flag.String("listen", "127.0.0.1:8000", "address to listen on")

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	flag.Parse()
	server, err := srv.New()
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}
	return server.Serve(*flagListenAddr)
}
