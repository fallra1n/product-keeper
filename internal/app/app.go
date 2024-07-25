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

	"github.com/fallra1n/product-keeper/config"
	"github.com/fallra1n/product-keeper/internal/adapters/authrepo"
	authrepoPg "github.com/fallra1n/product-keeper/internal/adapters/authrepo/postgres"
	"github.com/fallra1n/product-keeper/internal/adapters/productsrepo"
	productsrepoPg "github.com/fallra1n/product-keeper/internal/adapters/productsrepo/postgres"
	services "github.com/fallra1n/product-keeper/internal/core"
	"github.com/fallra1n/product-keeper/internal/core/auth"
	"github.com/fallra1n/product-keeper/internal/core/products"
	httphandler "github.com/fallra1n/product-keeper/internal/http-handler"
	"github.com/fallra1n/product-keeper/pkg/access"
	"github.com/fallra1n/product-keeper/pkg/logging"
	"github.com/fallra1n/product-keeper/pkg/postgresdb"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *sqlx.DB

	productsRepo products.ProductsRepo
	authRepo     auth.Authrepo

	httpServer *http.Server
}

func NewApp() *App {
	log.Println("env initializing")
	if err := godotenv.Load(); err != nil {
		log.Println("cannot loading .env file: " + err.Error())
	}

	cfg := config.MustLoad()

	a := &App{
		cfg:    cfg,
		logger: logging.SetupLogger(cfg.Env),
		db:     postgresdb.NewPostgresDB(access.PostgresConnect(cfg)),

		productsRepo: productsrepo.NewPostgresProducts(),
		authRepo:     authrepo.NewPostgresAuth(),
	}

	// tables init
	tx, err := a.db.Beginx()
	if err != nil {
		a.logger.Error("cannot start transaction:", err)
		os.Exit(1)
	}
	defer tx.Rollback()

	if err := productsrepoPg.CreateTable(tx); err != nil {
		a.logger.Error("cannot create products table in db:", err)
		os.Exit(1)
	}

	if err := authrepoPg.CreateTable(tx); err != nil {
		a.logger.Error("cannot create auth table in db:", err)
		os.Exit(1)
	}

	err = tx.Commit()
	if err != nil {
		a.logger.Error("cannot commit transaction:", err)
		os.Exit(1)
	}

	return a
}

func (a *App) Run() {
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

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.httpServer.Shutdown(ctx)
}
