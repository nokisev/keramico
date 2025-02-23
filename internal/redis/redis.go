package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	return &RedisClient{Client: rdb, ctx: ctx}
}

func (r *RedisClient) StoreToken(userID string, token string) error {
	return r.Client.Set(r.ctx, userID, token, 30*time.Minute).Err()
}

func (r *RedisClient) GetToken(userID string) (string, error) {
	return r.Client.Get(r.ctx, userID).Result()
}

func (r *RedisClient) DeleteToken(userID string) error {
	res := r.Client.Del(r.ctx, userID)
	if res != nil {
		return res.Err()
	}
	return nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
