package dbrepository

import (
	"context"
	"time"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/database"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/repository"
)

type mongoDbRepo struct {
	db  database.Database
	ctx context.Context
}

func NewMongoDbRepo(db database.Database, ctx context.Context) repository.Repository {
	return &mongoDbRepo{
		db:  db,
		ctx: ctx,
	}
}

const timeout = 3 * time.Second
