package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var DbName = os.Getenv("dbName")
var ColName = os.Getenv("collectionName")

func Migrate() {
	if DbName == "" || ColName == "" {
		log.Fatal("Service need specified 'dbName' and 'collectionName")
	}
	err := createCollection(DbName, ColName)
	if err != nil {
		switch err.(type) {
		case mongo.CommandError:
			log.Printf("Database '%s' and collection '%s' already exists", DbName, ColName)
		default:
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
