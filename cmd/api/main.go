package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/database"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/handlers"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/middlewares"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/routers"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	dsn  string
	port int
	addr string
)

var ctx context.Context = context.Background()

func init() {
	flag.StringVar(&dsn, "db", "mongodb://localhost:27017", "mongodb uri")
	flag.IntVar(&port, "p", 8000, "port to listen on")
	flag.Parse()
	addr = fmt.Sprintf(":%d", port)
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error in loading env files: %s\n", err.Error())
	}
}

func main() {
	handler, middleware, err := run()
	if err != nil {
		log.Fatal(err)
	}
	router := routers.NewRouter(handler, middleware)
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}
	log.Printf("Starting the server at %s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func run() (handlers.Handlers, middlewares.Middlewares, error) {
	database, err := database.NewDatabase(dsn, ctx)
	if err != nil {
		return nil, nil, err
	}
	redisPool := redis.NewClient(
		&redis.Options{
			Addr:         os.Getenv("REDIS_URL"),
			Password:     "",
			DB:           0,
			MaxIdleConns: 10,
			PoolSize:     10,
			MinIdleConns: 0,
		},
	)
	handler := handlers.NewHandlers(database, redisPool, ctx)
	middleware := middlewares.NewMiddlwares(redisPool, ctx)
	return handler, middleware, nil
}
