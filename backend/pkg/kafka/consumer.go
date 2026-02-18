package kafka

import (
	"context"
	"fmt"

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
		logger.SystemLogger.Info().
			Str("topic", message.Topic).
			Int32("partition", message.Partition).
			Int64("offset", message.Offset).
			Msg("Message claimed")

		if err := h.handler(message.Value); err != nil {
			logger.SystemLogger.Error().Err(err).Msg("Failed to process message")
			// Depending on requirement, we might not mark offset if processing fails
			// For now, we log and mark as processed to avoid infinite loops on bad messages
		}

		session.MarkMessage(message, "")
	}
	return nil
}
