package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func Get(dbName string, collectionName string, f bson.M, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	defer cancel()
	defer disconnect(ctx)

	return client.Database(dbName).Collection(collectionName).FindOne(context.TODO(), f, opts...), err
}

func GetCollection(dbName string, collectionName string) (*mongo.Collection, error) {
	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer disconnect(ctx)
	return client.Database(dbName).Collection(collectionName), err
}

func UpdateDocs(dbName string, collectionName string,
	f bson.D, u bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer disconnect(ctx)
	return client.Database(dbName).Collection(collectionName).UpdateMany(ctx, f, u, opts...)
}

func UpdateDoc(dbName string, collectionName string,
	f bson.D, u bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer disconnect(ctx)
	return client.Database(dbName).Collection(collectionName).UpdateOne(ctx, f, u, opts...)
}

func InsertDoc(dbName string, collectionName string,
	doc any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {

	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer disconnect(ctx)
	one, err := client.Database(dbName).Collection(collectionName).InsertOne(ctx, doc, opts...)
	return one, err
}

func DeleteDocs(dbName string, collectionName string, filter any, opts ...*options.DeleteOptions) (int64, error) {
	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	if err != nil {
		return 0, err
	}
	defer cancel()
	defer disconnect(ctx)
	many, err := client.Database(dbName).Collection(collectionName).DeleteMany(ctx, filter, opts...)
	return many.DeletedCount, err
}

func DeleteDoc(dbName string, collectionName string, filter any, opts ...*options.DeleteOptions) bool {
	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	defer cancel()
	defer disconnect(ctx)
	one, err := client.Database(dbName).Collection(collectionName).DeleteOne(ctx, filter, opts...)
	if err != nil {
		log.Println(err)
	}
	return one.DeletedCount == 1
}
