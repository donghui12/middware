package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, toValue interface{}) error
}
type RDB struct {
	client *redis.Client
}

var (
	rdb  *RDB
	once sync.Once
)

func GetClient() Client {
	return rdb
}

func Init(conf *Config) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DB,
	})
	pong, err := client.Ping(context.TODO()).Result()
	if err != nil {
		// TODO log.Errorf( pong )
		fmt.Println("this is result from ping:", pong, err)
		return fmt.Errorf("error to connect client")
	}
	once.Do(func() {
		rdb = &RDB{client: client}
	})
	return nil
}

func (r *RDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	valueJson, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := rdb.client.Set(ctx, key, valueJson, expiration).Err(); err != nil {
		return fmt.Errorf("set value to redis error %w", err)
	}
	return nil
}

func (r *RDB) Get(ctx context.Context, key string, toValue interface{}) error {
	userFromRedis, err := rdb.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("%s does not exist", key)
	}
	if err != nil {
		return fmt.Errorf("get value from redis error %w", errors.Unwrap(err))
	}
	if err = json.Unmarshal(userFromRedis, &toValue); err != nil {
		// TODO log err
		return errors.New("unmarshal redis data to struct error")
	}
	return nil
}
