package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/internal/domain"
)

type EventHandler struct {
	eventProcessor EventProcessor
}

func NewEventHandler(eventProcessor EventProcessor) *EventHandler {
	return &EventHandler{eventProcessor: eventProcessor}
}

func (h *EventHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return h.eventProcessor.Execute(ctx, domain.Event{
		Key:       msg.Key,
		Value:     msg.Value,
		Timestamp: msg.Timestamp,
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
	})
}
