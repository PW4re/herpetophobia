package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type disconnectFunc func(ctx context.Context)

var client *mongo.Client

func createClient() (*mongo.Client, disconnectFunc, error) {
	var err error
	f := func(ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}
	if client == nil {
		client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@mongo:27017/")) // stubbed // todo from env
	}
	return client, f, err
}

func connect(client *mongo.Client, er error) (ctx context.Context, cancel context.CancelFunc, err error) {
	if er != nil {
		return nil, nil, er
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	return
}
