package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
)

type ConsumerHandler func(message []byte) error

type Consumer interface {
	Consume(ctx context.Context, topic string, handler ConsumerHandler) error
	Close() error
}

type consumer struct {
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(cfg *config.Config, groupID string) (Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, groupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer group: %w", err)
	}

	return &consumer{consumerGroup: client}, nil
}

func (c *consumer) Consume(ctx context.Context, topic string, handler ConsumerHandler) error {
	consumerHandler := &saramaConsumerHandler{
		handler: handler,
	}

	for {
		if err := c.consumerGroup.Consume(ctx, []string{topic}, consumerHandler); err != nil {
			if err == sarama.ErrClosedConsumerGroup {
				return nil
			}
			logger.SystemLogger.Error().Err(err).Msg("Error from consumer")
			return err
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *consumer) Close() error {
	return c.consumerGroup.Close()
}

type saramaConsumerHandler struct {
	handler ConsumerHandler
}

func (h *saramaConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *saramaConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *saramaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		start := time.Now()
		
		var status int = 200
		var processErr error
		
		if err := h.handler(message.Value); err != nil {
			status = 500
			processErr = err
		}

		elapsed := time.Since(start)
		truncatedReq := logger.Truncate(string(message.Value), 65536)

		logger.SystemLogger.Info().
			Str("method", "KAFKA:CONSUME").
			Str("path", message.Topic).
			Int("status_code", status).
			Int64("latency", elapsed.Milliseconds()).
			Str("request_body", truncatedReq).
			Str("response_body", fmt.Sprintf("partition:%d, offset:%d", message.Partition, message.Offset)).
			Msg("kafka operation")

		if logger.SystemLogRepo != nil {
			_ = logger.SystemLogRepo.Create(session.Context(), &logger.SystemLog{
				Method:       "KAFKA:CONSUME",
				Path:         message.Topic,
				StatusCode:   status,
				Latency:      elapsed.Milliseconds(),
				RequestBody:  truncatedReq,
				ResponseBody: fmt.Sprintf("partition:%d, offset:%d", message.Partition, message.Offset),
			})
		}

		if processErr != nil {
			logger.SystemLogger.Error().Err(processErr).Msg("Failed to process message")
		}

		session.MarkMessage(message, "")
	}
	return nil
}
