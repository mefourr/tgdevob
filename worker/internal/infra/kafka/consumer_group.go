package kafka

import (
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/config"
)

func NewConsumerGroup(cfg config.Config) (sarama.ConsumerGroup, error) {
	consConfig := sarama.NewConfig()
	consConfig.Producer.Return.Errors = true
	consConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	consConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(cfg.Kafka.BootstrapServers, cfg.Kafka.Group, consConfig)
}
