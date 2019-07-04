package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

// BuildRedisOptions reads the configuration for the redis client
func BuildRedisOptions() *redis.Options {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		db = 0
	}
	opts := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       db,
	}
	return &opts
}

// BuildRedisOptions reads the configuration for the redis client
func BuildRedisClient(opts *redis.Options, check bool) *redis.Client {
	client := redis.NewClient(opts)
	if check {
		pong, err := client.Ping().Result()
		log.Println(pong, err)
	}
	return client
}
