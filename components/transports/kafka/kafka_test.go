package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func TestReader(t *testing.T) {
	cfg := Config{
		Brokers: []string{"localhost:9092"},
	}
	topic := "my-topic"
	go func() {
		writer := kafka.Writer{
			Addr: kafka.TCP(cfg.Brokers...),
		}
		for {
			time.Sleep(time.Second)
			err := writer.WriteMessages(context.Background(), kafka.Message{
				Topic: topic,
				Value: []byte("hello world"),
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	readerConfig := kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   topic,
	}

	reader := kafka.NewReader(readerConfig)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		fmt.Printf("Message at offset %d: %s\n", msg.Offset, string(msg.Value))
	}
}
