package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"

	"github.com/fallra1n/product-keeper/internal/core/products"
)

type ProductsStatistics struct {
	mq sarama.SyncProducer
}

func NewProducts(mq sarama.SyncProducer) *ProductsStatistics {
	return &ProductsStatistics{mq: mq}
}

func (s *ProductsStatistics) Send(p products.Product) error {
	pJson, err := json.Marshal(p)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic:     "products_statistics",
		Partition: -1,
		Value:     sarama.ByteEncoder(pJson),
	}

	if _, _, err := s.mq.SendMessage(&msg); err != nil {
		return err
	}

	return nil
}
