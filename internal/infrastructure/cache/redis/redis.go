package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/Ardnh/be-project-app/internal/config"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func NewRedisDB(cfg *config.Config) *redis.Client {

	hostPort := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     hostPort,           // contoh: "localhost:6379"
		Password: cfg.Redis.Password, // kosong kalau tidak pakai password
		DB:       cfg.Redis.DB,       // biasanya 0
	})

	// Test koneksi
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Redis connected successfully!")
	return rdb
}
