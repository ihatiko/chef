package kafka_consumer

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers                []string
	MaxWait                time.Duration
	WriteTimeOut           time.Duration
	ReadTimeOut            time.Duration
	HeartbeatInterval      time.Duration
	PartitionWatchInterval time.Duration
	DialTimeout            time.Duration
	AsyncDefault           bool
	MaxAttempts            int
	CommitIntervalDefault  int
	QueueCapacity          int
	MinBytes               float64 // 10KB
	MaxBytes               float64 // 10MB
	Compression            kafka.Compression
	Async                  bool
	UseSSL                 bool
	SslCaPem               string
	UseSASL                bool
	Username               string
	Password               string
	AllowAutoTopicCreation bool
	RequiredAcks           kafka.RequiredAcks
	BatchTimeoutMS         time.Duration
}
