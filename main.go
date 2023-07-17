package main

import (
	"fmt"

	kfk "github.com/segmentio/kafka-go"

	"middware/plugins/kafka"
	"middware/plugins/redis"
)

type Config struct {
	Redis redis.Config `yaml:"redis"`
	Kafka kafka.Config `yaml:"kafka"`
}

func initMid(config *Config) error {
	if err := redis.Init(&config.Redis); err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	if err := kafka.Init(&config.Kafka); err != nil {
		return fmt.Errorf("kafka: %w", err)
	}

	return nil
}

//redis.Config{
//	Host:     "113.24.60.157",
//	Port:     30706,
//	Username: "default",
//	Password: "Kjgnj93JKj3je",
//	DB:       1,
//}

func main() {
	// pass
	conf := Config{
		Kafka: kafka.Config{
			Brokers: nil,
			Topic:   "",
			Offset:  0,
		},
	}

	err := initMid(&conf)
	if err != nil {
		panic(err)
	}

	// redis
	//ctx := context.Background()
	//loginUser := entities.User{
	//	ID:        "1",
	//	Name:      "jdh",
	//	Age:       "24",
	//	LoginTime: "2022-12-22",
	//}
	//err = redis.GetClient().Set(ctx, loginUser.ID, loginUser, time.Hour)
	//if err != nil {
	//	panic(err)
	//}
	//
	//time.Sleep(time.Second * 2)
	//
	//user := entities.User{}
	//if err := redis.GetClient().Get(ctx, "1", &user); err != nil {
	//	panic(err)
	//}
	//fmt.Printf("login user: %+v", user)

	// kafka
	flag := make(chan bool, 1)
	messageChan := make(chan *kfk.Message, 10)

	go kafka.Client().Read(flag, messageChan)

	for message := range messageChan {
		fmt.Println(message)
	}
}
