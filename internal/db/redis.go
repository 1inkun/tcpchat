package db

import (
	"context"
	"log"

	"github.com/1inkun/tcpchat/config"
	"github.com/redis/go-redis/v9"
)

func InitRedisConn(cfg config.Config) (*redis.Client, error) {
	// TODO
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})
	ctx := context.Background()
	cmd := rdb.Ping(ctx)
	if cmd.Err() != nil {
		log.Printf("连接Redis数据库时出错:%s", cmd.Err())
		return nil, cmd.Err()
	}
	return rdb, nil
}
