package main

import (
	"log"
	"os"
	"syscall"

	"github.com/fallra1n/product-service/internal"
	"github.com/fallra1n/product-service/pkg/shutdown"

	"github.com/fallra1n/product-service/internal/config"
	"github.com/fallra1n/product-service/pkg/logging"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("env initializing")
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file: " + err.Error())
	}

	log.Println("config initializing")
	cfg := config.MustLoad()

	log.Println("logger initializing")
	logger := logging.SetupLogger(cfg.Env)

	logger.Info("creating application")
	app, err := internal.NewApp(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("running application")
	go app.Run()

	shutdown.Graceful([]os.Signal{syscall.SIGINT, syscall.SIGTERM}, app)
}
