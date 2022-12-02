package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoClient *mongo.Client
var PostsColl *mongo.Collection

func ConnectDb(uri string) {
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	MongoClient = client
}

func InitCollections(dbName string) {
	PostsColl = GetCollection(MongoClient, "posts", dbName)
}

func GetCollection(client *mongo.Client, colName string, dbName string) *mongo.Collection {
	return client.Database(dbName).Collection(colName)
}

func DisconnectDb() {
	func() {
		if err := MongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
