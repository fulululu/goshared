// Package redis ...
package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RawConfig struct {
	Host     string `envconfig:"REDIS_HOST,default=127.0.0.1"`
	Port     string `envconfig:"REDIS_PORT,default=6379"`
	User     string `envconfig:"REDIS_USER,default=default"`
	Password string `envconfig:"REDIS_PASSWORD,optional"`
	Database int    `envconfig:"REDIS_DATABASE,default=0"`
}

func NewClient(ctx context.Context, cfg RawConfig) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.Database,
	})
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return cli, err
	}

	return cli, nil
}
