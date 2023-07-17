package kafka

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"log"
	"sync"
)

type KFK struct {
	client *kafka.Reader
}

var (
	reader *KFK
	once   sync.Once
)

func Client() *KFK {
	return reader
}

func Init(conf *Config) error {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   conf.Brokers,
		Topic:     conf.Topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	err := reader.client.SetOffset(conf.Offset)
	if err != nil {
		return err
	}
	once.Do(func() {
		reader = &KFK{
			client: r,
		}
	})
	return nil
}

func (k *KFK) Read(flag chan bool, data chan *kafka.Message) {
	for {
		select {
		case <-flag:
			if err := k.client.Close(); err != nil {
				log.Fatal("failed to close reader:", err)
			}
			fmt.Println("break")
			return
		default:
		}
		m, err := reader.client.ReadMessage(context.Background())
		if err != nil {
			break
		}
		data <- &m
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
