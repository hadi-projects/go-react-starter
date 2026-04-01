package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hadi-projects/go-react-starter/config"
	"github.com/hadi-projects/go-react-starter/internal/worker"
	"github.com/hadi-projects/go-react-starter/pkg/kafka"
	"github.com/hadi-projects/go-react-starter/pkg/logger"
	"github.com/hadi-projects/go-react-starter/pkg/mailer"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(&cfg)

	logger.SystemLogger.Info().Msg("Starting Worker Service...")

	// Initialize Mailer
	mailService := mailer.NewMailer(&cfg)

	// Initialize Kafka Consumer
	// Group ID should be unique for this worker service
	consumer, err := kafka.NewConsumer(&cfg, "go-starter-worker-group")
	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Failed to create Kafka consumer")
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGTERM/SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.SystemLogger.Info().Msg("Shutting down worker...")
		cancel()
	}()

	topic := cfg.Kafka.Topic
	if topic == "" {
		topic = "password-reset"
	}

	logger.SystemLogger.Info().Str("topic", topic).Msg("Listening for messages...")

	err = consumer.Consume(ctx, topic, func(message []byte) error {
		return worker.ProcessResetPassword(message, &cfg, mailService)
	})

	if err != nil {
		logger.SystemLogger.Fatal().Err(err).Msg("Consumer error")
	}

	logger.SystemLogger.Info().Msg("Worker exited")
}
