package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/fallra1n/product-keeper/config"
	"github.com/fallra1n/product-keeper/internal/adapters/authrepo"
	"github.com/fallra1n/product-keeper/internal/adapters/productsrepo"
	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/internal/core/products"
	httphandler "github.com/fallra1n/product-keeper/internal/handler/http"
	authhttphandler "github.com/fallra1n/product-keeper/internal/handler/http/auth"
	productshttphandler "github.com/fallra1n/product-keeper/internal/handler/http/products"
	"github.com/fallra1n/product-keeper/pkg/access"
	"github.com/fallra1n/product-keeper/pkg/kafka"
	"github.com/fallra1n/product-keeper/pkg/logging"
	"github.com/fallra1n/product-keeper/pkg/postgresdb"
)

// App application
type App struct {
	cfg               *config.Config
	log               *slog.Logger
	db                *sqlx.DB
	kafkaSyncProducer sarama.SyncProducer

	authRepo     auth.AuthRepo
	productsRepo products.ProductsRepo

	authService     *auth.AuthService
	productsService *products.ProductsService

	authHandler     httphandler.AuthHandler
	productsHandler httphandler.ProductsHandler

	httpServer *http.Server
}

// NewApp creating new app
func NewApp() (*App, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("cannot loading .env file: %s", err)
	}

	cfg := config.MustLoad()

	a := &App{
		cfg:               cfg,
		log:               logging.SetupLogger(cfg.Env),
		db:                postgresdb.NewPostgresDB(access.PostgresConnect(cfg), cfg.Postgres.Timeout),
		kafkaSyncProducer: kafka.NewKafkaSyncProducer(access.KafkaConnect(cfg)),

		productsRepo: productsrepo.NewPostgresProducts(),
		authRepo:     authrepo.NewPostgresAuth(),
	}

	// services init
	a.authService = auth.NewAuthService(a.db, a.log, a.authRepo)
	a.productsService = products.NewProductsService(a.db, a.log, a.productsRepo)

	// http handlers init
	a.authHandler = authhttphandler.NewAuthHandler(a.log, a.authService)
	a.productsHandler = productshttphandler.NewProductsHandler(a.log, a.productsService)

	// http server init
	router := httphandler.SetupRouter(a.log, a.authHandler, a.productsHandler)

	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", a.cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
	}

	return a, nil
}

func (a *App) Run() {
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Error(fmt.Sprintf("error ocurred while running http-server server: %s", err))
		os.Exit(1)
	}
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.httpServer.Shutdown(ctx)
}
