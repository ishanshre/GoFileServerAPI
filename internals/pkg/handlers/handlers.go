package handlers

import (
	"context"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/database"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/repository"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/repository/dbrepository"
	"github.com/redis/go-redis/v9"
)

type Handlers interface{}

type handlers struct {
	mg          repository.Repository
	ctx         context.Context
	redisClient *redis.Client
}

func NewHandlers(database database.Database, r *redis.Client, ctx context.Context) Handlers {
	return &handlers{
		mg:          dbrepository.NewMongoDbRepo(database, ctx),
		redisClient: r,
		ctx:         ctx,
	}
}
