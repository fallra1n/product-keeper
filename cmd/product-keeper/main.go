package main

import (
	"os"
	"syscall"

	"github.com/fallra1n/product-keeper/internal/app"
	"github.com/fallra1n/product-keeper/pkg/shutdown"
)

func main() {
	appl, err := app.NewApp()
	if err != nil {
		os.Exit(1)
	}
	
	go appl.Run()
	shutdown.Graceful([]os.Signal{syscall.SIGINT, syscall.SIGTERM}, appl)
}
