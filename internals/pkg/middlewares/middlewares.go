package middlewares

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Middlewares interface{}

type middlewares struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewMiddlwares(r *redis.Client, ctx context.Context) Middlewares {
	return &middlewares{
		redisClient: r,
		ctx:         ctx,
	}
}
