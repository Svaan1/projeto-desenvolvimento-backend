package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNil = errors.New("no matching record found in redis database")
)

type RedisDB struct {
	Client *redis.Client
}

func (r *RedisDB) NewDatabase(ctx context.Context) (Database, error) {
	return NewRedisDB(ctx)
}

func NewRedisDB(ctx context.Context) (*RedisDB, error) {
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	log.Printf("Connecting to Redis at: redis:%s", redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("redis:%s", redisPort),
		Password: redisPassword,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisDB{Client: client}, nil
}
