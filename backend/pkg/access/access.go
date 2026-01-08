package access

import (
	"fmt"

	"github.com/fallra1n/product-keeper/config"
)

// PostgresConnect get connection string to postgres
func PostgresConnect(cfg *config.Config) string {
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	return connect
}

// PostgresTestConnect get connection string to postgres for testing
func PostgresTestConnect(cfg *config.Config) string {
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresTest.Host, cfg.PostgresTest.Port, cfg.PostgresTest.User, cfg.PostgresTest.Password, cfg.PostgresTest.DBName)

	return connect
}

// KafkaConnect get connection string to kafka
func KafkaConnect(cfg *config.Config) []string {
	urlList := make([]string, len(cfg.BrokerList))

	for id, row := range cfg.BrokerList {
		urlList[id] = fmt.Sprintf("%s:%s", row.Host, row.Port)
	}

	return urlList
}
