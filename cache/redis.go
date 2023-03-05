package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
)

// RedisCache is the redis cache configuration
type RedisCache struct {
	Host     string
	Password string
	DB       int
	Expires  time.Duration
}

// NewRedisCache creates a new redis cache
func NewRedisCache(config *RedisCache) *RedisCache {
	return &RedisCache{
		Host:     config.Host,
		Password: config.Password,
		DB:       config.DB,
		Expires:  config.Expires,
	}
}

// GetClient returns the redis client
func (cache *RedisCache) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.Host,
		Password: cache.Password,
		DB:       cache.DB,
	})
}

// GetValue gets a user from the cache
func (c *RedisCache) GetValue(key string) (interface{}, error) {
	client := c.GetClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var value interface{}
	err = json.Unmarshal([]byte(val), &value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// SetValue sets a user in the cache
func (c *RedisCache) SetValue(key string, value interface{}) error {
	client := c.GetClient()

	json, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	err = client.Set(key, json, c.Expires*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
