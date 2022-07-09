package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var Ctx = context.Background()

func CreateClient(int dbNo) *redis.client {
	rdb := redis.newClient(&redis.Options{
		Addr : os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB: dbNo,
	})
	return rdb

}