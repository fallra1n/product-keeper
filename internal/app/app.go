package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fallra1n/product-service/internal/config"
	httpServer "github.com/fallra1n/product-service/internal/http-server"
	"github.com/fallra1n/product-service/internal/http-server/handlers"
)

type App interface {
	Run()
	Close() error
}

type app struct {
	cfg        *config.Config
	logger     *slog.Logger
	httpServer *http.Server
}

func NewApp(cfg *config.Config, logger *slog.Logger) App {
	return &app{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *app) Run() {
	ah := handlers.NewAuthHandler()
	prh := handlers.NewProductHandler()

	router := httpServer.SetupRouter(ah, prh, a.logger)

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(a.cfg.HTTPServer.Address),
		Handler: router,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("error ocurred while running http-server server: %s", err.Error())
			os.Exit(1)
		}
	}()
}

func (a *app) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.httpServer.Shutdown(ctx)
}
