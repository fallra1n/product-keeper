package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// NewSyncProducer get new kafka sync producer
func NewSyncProducer(urlList []string) sarama.SyncProducer {
	cfg := sarama.NewConfig()

	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(urlList, cfg)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %s", err)
	}

	return conn
}

// NewAsyncProducer get new kafka async producer
func NewAsyncProducer(urlList []string) sarama.AsyncProducer {
	cfg := sarama.NewConfig()

	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5

	conn, err := sarama.NewAsyncProducer(urlList, cfg)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %s", err)
	}

	return conn
}
