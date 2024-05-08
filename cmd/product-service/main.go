package main

import (
	"log"
	"os"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/fallra1n/product-service/internal/app"
	"github.com/fallra1n/product-service/internal/config"
	"github.com/fallra1n/product-service/pkg/logging"
	"github.com/fallra1n/product-service/pkg/shutdown"
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
	appl := app.NewApp(cfg, logger)

	logger.Info("running application")
	go appl.Run()

	shutdown.Graceful([]os.Signal{syscall.SIGINT, syscall.SIGTERM}, appl)
}
