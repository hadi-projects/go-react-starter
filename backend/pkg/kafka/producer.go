package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type Producer interface {
	Publish(topic string, message interface{}) error
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

func (p *producer) Publish(topic string, message interface{}) error {
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

	logger.SystemLogger.Info().
		Str("topic", topic).
		Int32("partition", partition).
		Int64("offset", offset).
		Msg("Message published to Kafka")

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
