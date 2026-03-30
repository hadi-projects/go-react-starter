package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type Producer interface {
	Publish(ctx context.Context, topic string, message interface{}) error
	Close() error
	Status() string
}

type producer struct {
	syncProducer sarama.SyncProducer
}

func NewProducer(cfg *config.Config) (Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	syncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &producer{syncProducer: syncProducer}, nil
}

func (p *producer) Publish(ctx context.Context, topic string, message interface{}) error {
	start := time.Now()
	requestID := ""
	if rid, ok := ctx.Value(logger.CtxKeyRequestID).(string); ok {
		requestID = rid
	}
	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	elapsed := time.Since(start)
	truncatedReq := logger.Truncate(string(value), 65536)

	logger.SystemLogger.Info().
		Str("method", "KAFKA:PUBLISH").
		Str("path", topic).
		Int("status_code", 200).
		Int64("latency", elapsed.Milliseconds()).
		Str("request_body", truncatedReq).
		Str("response_body", fmt.Sprintf("partition:%d, offset:%d", partition, offset)).
		Msg("kafka operation")

	if logger.SystemLogRepo != nil {
		_ = logger.SystemLogRepo.Create(ctx, &logger.SystemLog{
			RequestID:    requestID,
			Method:       "KAFKA:PUBLISH",
			Path:         topic,
			StatusCode:   200,
			Latency:      elapsed.Milliseconds(),
			RequestBody:  truncatedReq,
			ResponseBody: fmt.Sprintf("partition:%d, offset:%d", partition, offset),
		})
	}

	return nil
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}

func (p *producer) Status() string {
	if p.syncProducer == nil {
		return "disconnected"
	}
	return "connected"
}
