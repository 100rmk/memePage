package mongoDB

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func GetPostsByRating(limit int64, date primitive.DateTime) []bson.M {
	opts := options.Find().SetLimit(limit).SetSort(bson.M{"likes_count": -1})
	filter := bson.D{
		{"timestamp", bson.D{{"$gte", date}}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cursor, err := PostsRatingColl.Find(ctx, filter, opts)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	var res []bson.M
	if err = cursor.All(ctx, &res); err != nil {
		log.Println(err)
		panic(err)
	}

	return res
}
