package client

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sample/config"
)

type StorageType struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
	ctx         context.Context
}

func CreateStorage(cfg *config.AppConfig, ctx context.Context) *StorageType {
	mongoClient := mongoConnect(cfg.MongoUri, ctx)
	collection := mongoClient.Database(cfg.MongoDbName).Collection(cfg.MongoCollectionName)
	return &StorageType{
		mongoClient: mongoClient,
		collection:  collection,
		ctx:         ctx,
	}
}

func Add(storage *StorageType, document interface{}) error {
	_, err := storage.collection.InsertOne(storage.ctx, document)
	return err
}

func FindOne(storage *StorageType, filter interface{}, resultType interface{}) error {
	return storage.collection.FindOne(storage.ctx, filter).Decode(resultType)
}

func mongoConnect(mongoUri string, ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}
	return client
}
