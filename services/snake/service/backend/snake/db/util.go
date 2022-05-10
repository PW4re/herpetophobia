package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type disconnectFunc func(ctx context.Context)

func createClient() (*mongo.Client, disconnectFunc, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://admin:admin@mongo:27017/")) // stubbed

	return client, func(ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}, err
}

func connect(client *mongo.Client, er error) (ctx context.Context, cancel context.CancelFunc, err error) {
	if er != nil {
		return nil, nil, er
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	return
}
