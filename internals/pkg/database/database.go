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
	Close()
	StartTranscation() error
	EndTranscation() error
	CommitTranscation() error
}

type database struct {
	client  *mongo.Client
	session mongo.Session
	ctx     context.Context
}

func NewDatabase(dsn string, ctx context.Context) (Database, error) {
	client, session, err := Connect(dsn, ctx)
	if err != nil {
		return nil, err
	}
	return &database{
		client:  client,
		session: session,
		ctx:     ctx,
	}, nil
}

func Connect(dsn string, ctx context.Context) (*mongo.Client, mongo.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	session, err := client.StartSession()
	if err != nil {
		return nil, nil, err
	}
	return client, session, nil
}

func (d *database) Close() {
	if d.session != nil {
		d.session.EndSession(d.ctx)
	}
	d.client.Disconnect(d.ctx)
}

func (d *database) GetUserCollection() *mongo.Collection {
	return d.client.Database("myDB").Collection("users")
}

func (d *database) GetFileCollection() *mongo.Collection {
	return d.client.Database("myDB").Collection("files")
}

func (d *database) StartTranscation() error {
	return d.session.StartTransaction()
}

func (d *database) EndTranscation() error {
	return d.session.AbortTransaction(d.ctx)
}

func (d *database) CommitTranscation() error {
	return d.session.CommitTransaction(d.ctx)
}
