package kafka_producer

import (
	"github.com/ihatiko/chef/components/clients/kafka"
)

type Config struct {
	kafka.Config
	Topic string
}
