package main

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Logger struct {
	topic    string
	producer *kafka.Producer
}

func (logger *Logger) Initialize(kafkaAddr, kafkaTopic string) (err error) {
	logger.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaAddr,
	})
	if err != nil {
		return err
	}
	logger.topic = kafkaTopic
	return nil
}

func (logger *Logger) Close() error {
	logger.producer.Close()
	return nil
}

func (logger *Logger) Log(msg string) {
	_ = logger.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &logger.topic},
		Value:          []byte(msg),
	}, nil)
}
