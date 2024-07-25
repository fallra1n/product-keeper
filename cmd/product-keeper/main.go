package main

import (
	"os"
	"syscall"

	"github.com/fallra1n/product-keeper/internal/app"
	"github.com/fallra1n/product-keeper/pkg/shutdown"
)

func main() {
	appl := app.NewApp()
	go appl.Run()
	shutdown.Graceful([]os.Signal{syscall.SIGINT, syscall.SIGTERM}, appl)
}
