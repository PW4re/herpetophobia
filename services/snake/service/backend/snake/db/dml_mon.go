package db

import (
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

	return client.Database(dbName).Collection(collectionName).FindOne(ctx, f, opts...), err
}

func GetResList(dbName string, collectionName string, f bson.D, opts ...*options.FindOptions) ([]bson.D, error) {
	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer disconnect(ctx)
	cur, err := client.Database(dbName).Collection(collectionName).Find(ctx, f, opts...)
	var results []bson.D
	err = cur.All(ctx, results)
	return results, err
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
