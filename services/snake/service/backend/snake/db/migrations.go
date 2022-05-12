package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var DbName = os.Getenv("dbName")
var ColName = os.Getenv("collectionName")

func Migrate() {
	//TODO: get from env
	err := createCollection(DbName, ColName)
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
	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	//if err != nil {
	//	return err
	//}
	defer cancel()
	defer disconnect(ctx)
	client.Database(dbName).Collection(collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{},
		Options: options.Index().SetExpireAfterSeconds(15.5 * 60)})
}
