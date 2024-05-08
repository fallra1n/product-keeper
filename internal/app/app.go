package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fallra1n/product-service/internal/config"
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

func NewApp(cfg *config.Config, logger *slog.Logger) (App, error) {
	return &app{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (a *app) Run() {
	router := gin.Default()
	router.GET("/home", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(a.cfg.HTTPServer.Address),
		Handler: router,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("error ocurred while running http server: %s", err.Error())
			os.Exit(1)
		}
	}()
}

func (a *app) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.httpServer.Shutdown(ctx)
}
