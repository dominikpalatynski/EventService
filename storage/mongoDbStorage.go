package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dominikpalatynski/EventService/types"
	"github.com/dominikpalatynski/EventService/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbStorage struct {
	db *MonogCollections
}

type MonogCollections struct {
	Events *mongo.Collection
	Contents *mongo.Collection
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
	Contentcollection := db.Collection("contents")

	return &MongoDbStorage{
		db: &MonogCollections{
			Events: Eventcollection,
			Contents: Contentcollection,
		},
	}, nil
}

func (s *MongoDbStorage) GetEvents(matcher map[string]interface{}) ([]types.Event, error){
	var events []types.Event
	filter := bson.M{
		"user_id": matcher["user_id"],
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

		var rawEvent bson.M

        if err := cursor.Decode(&rawEvent); err != nil {
            fmt.Printf("error here 2: %v", err.Error())
            return nil, err
        }

        fmt.Printf("rawEvent: %+v\n", rawEvent) // Logowanie surowego dokumentu

		var event types.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
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

func (s *MongoDbStorage) AddEvent(event *types.Event) (error) {

	initializedContent := &types.Content{Homework: "", Note: ""}

	insertedContent, err := s.db.Contents.InsertOne(context.Background(), initializedContent)

	if err != nil {
		return err
	}

	event.ContentId = insertedContent.InsertedID.(primitive.ObjectID)
	
	insertResult, err := s.db.Events.InsertOne(context.Background(), event)

	if err != nil {
		return err
	}

	event.ID = insertResult.InsertedID.(primitive.ObjectID)

	return nil
}

func (s *MongoDbStorage) UpdateEventById(eventId primitive.ObjectID, updatedData map[string]interface{}) (error) {

	filter := bson.M{"_id":eventId}
	update := bson.M{"$set": updatedData}

	if _, err :=  s.db.Events.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}

	return nil
}

func (s *MongoDbStorage) DeleteEventById(eventId primitive.ObjectID) (error) {
	filter := bson.M{"_id":eventId}

	if _, err :=  s.db.Events.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return nil
}

func (s *MongoDbStorage) GetContents() ([]types.Content, error) {
	var contents []types.Content
	cursor, err := s.db.Contents.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()){
		var content types.Content
		if err := cursor.Decode(&content); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	return contents, nil
}

func (s *MongoDbStorage) UpdateContentById(contentId primitive.ObjectID, updatedData map[string]interface{}) (error) {

	filter := bson.M{"_id":contentId}
	update := bson.M{"$set": updatedData}

	if _, err :=  s.db.Contents.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}

	return nil
}
