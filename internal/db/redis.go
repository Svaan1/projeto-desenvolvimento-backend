package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	Client *redis.Client
}

func NewRedisDB(ctx context.Context) (*RedisDB, error) {
	redisUser := os.Getenv("REDIS_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	log.Printf("Connecting to Redis at: %s:%s", redisHost, redisPort)

	opt, err := redis.ParseURL(
		fmt.Sprintf("redis://%s:%s@%s:%s", redisUser, redisPassword, redisHost, redisPort),
	)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisDB{Client: client}, nil
}

// SetObject stores an object as a JSON string in Redis
func (r *RedisDB) SetObject(ctx context.Context, key string, obj interface{}) error {
	// Marshal the object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err // JSON marshaling error
	}

	// Set the JSON string in Redis
	err = r.Client.Set(ctx, key, data, 0).Err() // 0 means no expiration
	if err != nil {
		return err // Redis error
	}

	return nil
}

// GetObject retrieves a whole object from Redis by its key
func (r *RedisDB) GetObject(ctx context.Context, key string, obj interface{}) error {
	// Get the JSON string from Redis
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // Key does not exist
		}
		return err // Other error
	}

	// Unmarshal the JSON string into the provided obj
	err = json.Unmarshal([]byte(val), obj)
	if err != nil {
		return err // Unmarshal error
	}

	return nil // Successfully retrieved
}
