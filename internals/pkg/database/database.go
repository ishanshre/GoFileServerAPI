package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	GetUserCollection() *mongo.Collection
	GetFileCollection() *mongo.Collection
}

type database struct {
	client *mongo.Client
	ctx    context.Context
}

func NewDatabase(dsn string, ctx context.Context) (Database, error) {
	client, err := Connect(dsn, ctx)
	if err != nil {
		return nil, err
	}
	return &database{
		client: client,
		ctx:    ctx,
	}, nil
}

func Connect(dsn string, ctx context.Context) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

func (d *database) GetUserCollection() *mongo.Collection {
	return d.client.Database("myDB").Collection("users")
}

func (d *database) GetFileCollection() *mongo.Collection {
	return d.client.Database("myDB").Collection("files")
}
