package dml

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type disconnectFunc func(ctx context.Context)

func createClient() (*mongo.Client, disconnectFunc) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")) // stubbed
	if err != nil {
		log.Fatal(err)
	}

	return client, func(ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func connect(client *mongo.Client) (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Get(dbName string, collectionName string, f bson.M, opts ...*options.FindOneOptions) *mongo.SingleResult {
	client, disconnect := createClient()
	ctx, cancel := connect(client)
	defer cancel()
	defer disconnect(ctx)

	return client.Database(dbName).Collection(collectionName).FindOne(context.TODO(), f, opts...)
}

func CreateCollection(dbName string, name string, opts ...*options.CreateCollectionOptions) {
	client, disconnect := createClient()
	ctx, cancel := connect(client)
	defer cancel()
	defer disconnect(ctx)
	err := client.Database(dbName).CreateCollection(ctx, name, opts...)
	if err != nil {
		log.Fatal(err)
	}
}

func GetCollection(dbName string, collectionName string) *mongo.Collection {
	client, disconnect := createClient()
	ctx, cancel := connect(client)
	defer cancel()
	defer disconnect(ctx)
	return client.Database(dbName).Collection(collectionName)
}

func UpdateDocs(dbName string, collectionName string,
	f bson.D, u bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	client, disconnect := createClient()
	ctx, cancel := connect(client)
	defer cancel()
	defer disconnect(ctx)
	return client.Database(dbName).Collection(collectionName).UpdateMany(ctx, f, u, opts...)
}

func UpdateDoc(dbName string, collectionName string,
	f bson.D, u bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	client, disconnect := createClient()
	ctx, cancel := connect(client)
	defer cancel()
	defer disconnect(ctx)
	return client.Database(dbName).Collection(collectionName).UpdateOne(ctx, f, u, opts...)
}

func InsertDoc(dbName string, collectionName string, doc any, opts ...*options.InsertOneOptions) *mongo.InsertOneResult {
	client, disconnect := createClient()
	ctx, cancel := connect(client)
	defer cancel()
	defer disconnect(ctx)
	one, err := client.Database(dbName).Collection(collectionName).InsertOne(ctx, doc, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return one
}
