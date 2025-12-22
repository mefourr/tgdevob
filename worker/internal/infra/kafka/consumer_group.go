package kafka

import (
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/config"
)

func NewConsumerGroup(cfg config.Config) (sarama.ConsumerGroup, error) {
	consConf := sarama.NewConfig()
	consConf.Producer.Return.Errors = true
	consConf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	consConf.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(cfg.Kafka.BootstrapServers, cfg.Kafka.Group, consConf)
}
