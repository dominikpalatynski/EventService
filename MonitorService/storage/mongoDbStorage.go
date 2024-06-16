package storage

import (
	"context"
	"os"
	"time"

	"MonitorService/types"
	"MonitorService/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbStorage struct {
	db *MonogCollections
}

type MonogCollections struct {
	Events *mongo.Collection
}

func NewMongoDbStorage() (*MongoDbStorage, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	util.LoadEnv()

	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")),
	   )

	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	db := mongoClient.Database("EventServiceDb")

	Eventcollection := db.Collection("events")

	return &MongoDbStorage{
		db: &MonogCollections{
			Events: Eventcollection,
		},
	}, nil
}

func (s *MongoDbStorage) GetAllEvents(matcher map[string]interface{}) ([]types.Event, error){
	var events []types.Event
	filter := bson.M{
		"start_date": bson.M{
			"$gte": matcher["start_date"],
			"$lt": matcher["end_date"],
		},
	}
	
	cursor, err := s.db.Events.Find(context.Background(), filter)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()){
		var event types.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
