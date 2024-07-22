package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fallra1n/product-keeper/internal/config"
	httphandler "github.com/fallra1n/product-keeper/internal/http-handler"
	"github.com/fallra1n/product-keeper/internal/services"
	"github.com/fallra1n/product-keeper/internal/storage/postgres"
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
	s, err := postgres.New(a.cfg)
	if err != nil {
		a.logger.Error("failed to connecting to the database: " + err.Error())
		os.Exit(1)
	}

	if err := s.CreateTables(); err != nil {
		a.logger.Error("failed to create tables: " + err.Error())
		os.Exit(1)
	}

	servs := services.NewServices(s)

	ah := httphandler.NewAuthHandler(servs, a.logger)
	prh := httphandler.NewProductHandler(servs, a.logger)

	router := httphandler.SetupRouter(ah, prh, a.logger)

	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", a.cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
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
