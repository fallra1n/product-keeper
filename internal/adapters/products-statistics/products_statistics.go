package productsstatistics

import (
	"github.com/IBM/sarama"

	"github.com/fallra1n/product-keeper/internal/adapters/products-statistics/kafka"
)

func NewKafkaProducts(mq sarama.SyncProducer) *kafka.ProductsStatistics {
	return kafka.NewProducts(mq)
}
