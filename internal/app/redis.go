package app

import (
	"fmt"
	"time"
)

// RedisClient is a mock Redis client
type RedisClient struct {
	host     string
	port     int
	password string
	db       int
	enabled  bool
}

// Redis is the global Redis client
var Redis *RedisClient

// InitRedis initializes the Redis connection
func InitRedis() {
	// Default Redis configuration
	host := ConfigData.Redis.Host
	port := ConfigData.Redis.Port
	password := ConfigData.Redis.Password
	db := ConfigData.Redis.DB

	Redis = &RedisClient{
		host:     host,
		port:     port,
		password: password,
		db:       db,
		enabled:  true,
	}

	// In a real application, you would connect to Redis here
	// For now, just log that Redis is "connected"
	fmt.Printf("Redis initialized successfully (host=%s, port=%d)\n", host, port)
	Info("Redis initialized successfully", "host", host, "port", port)
}

// Set stores a key-value pair in Redis with an optional expiration time
func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
	if !r.enabled {
		return fmt.Errorf("redis is not enabled")
	}

	fmt.Printf("MOCK REDIS: Set %s=%s with expiration %v\n", key, value, expiration)
	return nil
}

// Get retrieves a value from Redis by key
func (r *RedisClient) Get(key string) (string, error) {
	if !r.enabled {
		return "", fmt.Errorf("redis is not enabled")
	}

	fmt.Printf("MOCK REDIS: Get %s\n", key)
	return "mock-value-for-" + key, nil
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
	if !r.enabled {
		return fmt.Errorf("redis is not enabled")
	}

	fmt.Printf("MOCK REDIS: Delete %s\n", key)
	return nil
}

// IsEnabled returns whether Redis is enabled
func (r *RedisClient) IsEnabled() bool {
	return r != nil && r.enabled
}

// GetTTL gets the TTL of a key
func (r *RedisClient) GetTTL(key string) (time.Duration, error) {
	if !r.enabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	fmt.Printf("MOCK REDIS: GetTTL %s\n", key)
	return time.Hour, nil // Mock 1 hour TTL
}

// Incr increments the integer value of a key by one
func (r *RedisClient) Incr(key string) (int64, error) {
	if !r.enabled {
		return 0, fmt.Errorf("redis is not enabled")
	}

	fmt.Printf("MOCK REDIS: Incr %s\n", key)
	return 1, nil // Mock value after increment
}

// Close closes the Redis client connection
func (r *RedisClient) Close() error {
	if !r.enabled {
		return fmt.Errorf("redis is not enabled")
	}

	fmt.Println("MOCK REDIS: Connection closed")
	return nil
}
