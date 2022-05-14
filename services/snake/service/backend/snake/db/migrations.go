package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func Migrate() {
	//TODO: get from env
	err := createCollection("snake", "level")
	if err != nil {
		_, ok := err.(mongo.CommandError)
		log.Println(err.Error())
		if ok {
			log.Printf("Database '%s' and collection '%s' already exists", "local", "test")
		} else {
			log.Fatal(err)
		}
	}
}

func createIndex(dbName string, collectionName string) {
	// TODO: разобраться, какие нужны индексы (точно нужен ttl-index)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	getCollection(dbName, collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{},
		Options: options.Index().SetExpireAfterSeconds(15.5 * 60)})
	cancel()
}
