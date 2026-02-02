package kafka_consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/mefourr/tgdevob/worker/internal/dto"
	"github.com/mefourr/tgdevob/worker/internal/processor/usecase"
	"github.com/mefourr/tgdevob/worker/pkg/logger"
	"log/slog"
)

type EventHandler struct {
	eventProcessor usecase.EventProcessor
}

func NewEventHandler(eventProcessor usecase.EventProcessor) *EventHandler {
	return &EventHandler{eventProcessor: eventProcessor}
}

func (h *EventHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	m, err := parse(msg.Value)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to unmarshal eventProcessor")
		return err
	}
	return h.eventProcessor.Execute(ctx, m)
}

var ErrNoData = errors.New("no data in consumed message")

func parse(data []byte) (*dto.Message, error) {
	if len(data) == 0 {
		return nil, ErrNoData
	}
	var res dto.Message
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
