package mongoDB

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var MClient *mongo.Client
var PostsRatingColl *mongo.Collection

func ConnectDb(uri string) {
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	MClient = client
}

func InitCollections(dbName string) {
	PostsRatingColl = GetCollection(MClient, "rating", dbName)
}

func GetCollection(client *mongo.Client, colName string, dbName string) *mongo.Collection {
	return client.Database(dbName).Collection(colName)
}

func CreateViews(dbName string) {
	db := MClient.Database(dbName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if !isCollectionExists(db, "rating") {
		pipeline := []bson.M{
			{"$project": bson.M{
				"file_id":        1,
				"timestamp":      1,
				"content_type":   1,
				"likes_count":    bson.M{"$size": "$likes"},
				"dislikes_count": bson.M{"$size": "$dislikes"},
			}},
		}
		err := db.CreateView(ctx, "rating", "posts", pipeline)

		if err != nil {
			panic(err)
		}
	}
}

func DisconnectDb() {
	func() {
		if err := MClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func isCollectionExists(db *mongo.Database, name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	for _, colName := range collections {
		if colName == name {
			return true
		}
	}

	return false
}
