package config

import (
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers                            string `env:"KAFKA_BROKERS,required"`
	TopicConfirmedReservation          string `env:"KAFKA_TOPIC_CONFIRMED_RESERVATION,required"`
	GroupInventoryConfirmedReservation string `env:"KAFKA_GROUP_INVENTORY_CONFIRMED_RESERVATION,required"`
}

func NewKafkaReader(cfg KafkaConfig, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(cfg.Brokers, ","),
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}
