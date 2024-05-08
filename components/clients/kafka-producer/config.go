package kafka_producer

import (
	"github.com/ihatiko/olymp/components/clients/kafka"
)

type Config struct {
	kafka.Config
	Topic string
}
