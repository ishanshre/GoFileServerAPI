package handlers

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/database"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/repository"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/repository/dbrepository"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/validators"
	"github.com/redis/go-redis/v9"
)

type Handlers interface {
	//admin interface
	GetUsers(w http.ResponseWriter, r *http.Request)
	AdminDeleteUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)

	//user interface
	UserRegister(w http.ResponseWriter, r *http.Request)
	UserLogin(w http.ResponseWriter, r *http.Request)
	UserLogout(w http.ResponseWriter, r *http.Request)
	GetMe(w http.ResponseWriter, r *http.Request)
	DeleteMe(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	mg          repository.Repository
	ctx         context.Context
	redisClient *redis.Client
}

var validate *validator.Validate

func NewHandlers(database database.Database, r *redis.Client, ctx context.Context) Handlers {
	validate = validator.New()
	validate.RegisterValidation("uppercase", validators.Uppercase)
	validate.RegisterValidation("lowercase", validators.LowerCase)
	validate.RegisterValidation("number", validators.Number)
	return &handlers{
		mg:          dbrepository.NewMongoDbRepo(database, ctx),
		redisClient: r,
		ctx:         ctx,
	}
}
