package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
	Ctx   = context.Background()
)

func ConnectRedis() {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = 0
	}

	Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	if err := Redis.Ping(Ctx).Err(); err != nil {
		log.Fatal("Failed to connect Redis:", err)
	}

	log.Println("Redis connected")
}
