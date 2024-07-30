package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// NewKafkaSyncProducer get new kafka producer
func NewKafkaSyncProducer(urlList []string) sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(urlList, cfg)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %s", err)
	}

	return conn
}
