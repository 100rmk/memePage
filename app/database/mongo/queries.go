package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func GetPostsByRating(limit int, results chan []bson.M) {
	pipeline := []bson.M{
		{"$project": bson.M{"file_id": 1, "timestamp": 1, "count": bson.M{"$size": "$likes"}}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": limit},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cursor, err := PostsColl.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	var res []bson.M
	if err = cursor.All(ctx, &res); err != nil {
		log.Println(err)
		panic(err)
	}

	results <- res
}
