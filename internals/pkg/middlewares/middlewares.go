package middlewares

import (
	"context"
	"net/http"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/repository"
	"github.com/redis/go-redis/v9"
)

type Middlewares interface {
	JwtAuth(next http.Handler) http.Handler
	CheckAdmin(next http.Handler) http.Handler
	FileAuth(next http.Handler) http.Handler
	FileOwner(next http.Handler) http.Handler
}

type middlewares struct {
	redisClient *redis.Client
	repository  repository.Repository
	ctx         context.Context
}

func NewMiddlwares(r *redis.Client, repository repository.Repository, ctx context.Context) Middlewares {
	return &middlewares{
		redisClient: r,
		repository:  repository,
		ctx:         ctx,
	}
}
