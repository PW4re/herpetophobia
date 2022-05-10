package db

import "go.mongodb.org/mongo-driver/mongo/options"

func createCollection(dbName string, name string, opts ...*options.CreateCollectionOptions) error {
	client, disconnect, err := createClient()
	ctx, cancel, err := connect(client, err)
	if err != nil {
		return err
	}
	defer cancel()
	defer disconnect(ctx)
	err = client.Database(dbName).CreateCollection(ctx, name, opts...)
	return err
}
