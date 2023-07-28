package middlewares

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Middlewares interface {
	JwtAuth(next http.Handler) http.Handler
	CheckAdmin(next http.Handler) http.Handler
}

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
