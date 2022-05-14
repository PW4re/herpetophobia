package connhelper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

var client *mongo.Client
var once sync.Once

func GetClient() *mongo.Client {
	var clientError error
	once.Do(func() {
		client, clientError = mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@127.0.0.1:27017/")) // stubbed
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := client.Connect(ctx) // significant panic reason
		if err != nil {
			clientError = err
		}
		err = client.Ping(ctx, nil)
		if err != nil {
			clientError = err
		}
		cancel()
	})
	if clientError != nil {
		log.Fatal(clientError)
	}
	return client
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Disconnect(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
}
