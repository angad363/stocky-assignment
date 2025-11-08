package price

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx       = context.Background()
	RedisConn *redis.Client
)

// InitRedis initializes the Redis client connection
func InitRedis() {
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // example: "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // empty if no password
		DB:       0,
	})

	_, err := RedisConn.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis successfully!")
}