package kafka

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/goccy/go-json"
)

const (
	PollTimeOutMs = 100
)

// Consumer is the generic consumer for kafka messages
type Consumer struct {
	KafkaConsumer *kafka.Consumer
	logger        logger
	config        ConsumerConfig
}

// New creates new Kafka consumer
func New(config ConsumerConfig, logger logger) *Consumer {
	return &Consumer{
		logger: logger,
		config: config,
	}
}

// Init uses for initializing Consumer
func (c *Consumer) Init() error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.config.Address,
		"group.id":          c.config.GroupID,
		"auto.offset.reset": c.config.AutoOffsetReset,
	})
	if err != nil {
		return fmt.Errorf("Kafka.Consumer.Init: failed to create: %s", err.Error())
	}

	if err := consumer.SubscribeTopics([]string{c.config.Topic}, nil); err != nil {
		return fmt.Errorf("Kafka.Consumer.Init: failed subscribe topic: %s because of: %s", c.config.Topic, err.Error())
	}

	c.KafkaConsumer = consumer

	return nil
}

// RunConsumer starts executing Kafka Consumer that sends deserialized response into output chan
func RunConsumer[TValue any](consumer *Consumer, logger logger, output chan TValue) error {
	defer func() {
		if err := consumer.KafkaConsumer.Close(); err != nil {
			logger.Errorf("Kafka.Consumer.RunConsumer: failed to close consumer: %s", err.Error())
		}
	}()

	for {
		event := consumer.KafkaConsumer.Poll(PollTimeOutMs)
		switch eventType := event.(type) {
		case *kafka.Message:
			logger.Infof("Kafka.Consumer.RunConsumer: received message: %s", eventType.String())

			if _, err := consumer.KafkaConsumer.Commit(); err != nil {
				logger.Infof("Kafka.Consumer.RunConsumer: commit for: %s failed because of: %s", err.Error(), eventType.String())
				continue
			}

			var object TValue
			if err := json.Unmarshal(eventType.Value, &object); err != nil {
				logger.Infof("Kafka.Consumer.RunConsumer: unmarshal for: %s failed because of: %s", err.Error(), eventType.String())
			}

			output <- object
		case kafka.PartitionEOF:
			logger.Infof("Kafka.Consumer.RunConsumer: partition end for: %s", eventType.String())
		case kafka.Error:
			logger.Infof("Kafka.Consumer.RunConsumer: failed to consume message: %s", eventType.String())
			return errors.New(eventType.Error())
		}
	}
}
