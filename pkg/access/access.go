package access

import (
	"fmt"

	"github.com/fallra1n/product-keeper/config"
)

func PostgresConnect(cfg *config.Config) string {
	connect := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	return connect
}
