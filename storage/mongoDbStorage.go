package storage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/dominikpalatynski/EventService/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Event struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId string `json:"user_id" bson:"user_id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Content json.RawMessage `json:"content,omitempty" bson:"content,omitempty"`
}

type MongoDbStorage struct {
	db *mongo.Collection
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

	collection := db.Collection("events")

	return &MongoDbStorage{
		db: collection,
	}, nil
}

func (s *MongoDbStorage) GetEvents() ([]Event, error){
	var events []Event
	cursor, err := s.db.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()){
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *MongoDbStorage) AddEvent(event *Event) (error) {

	insertResult, err := s.db.InsertOne(context.Background(), event)

	if err != nil {
		return err
	}

	event.ID = insertResult.InsertedID.(primitive.ObjectID)

	return nil
}

func (s *MongoDbStorage) UpdateById(eventId primitive.ObjectID, updatedData map[string]interface{}) (error) {

	filter := bson.M{"_id":eventId}
	update := bson.M{"$set": updatedData}

	if _, err := s.db.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}

	return nil
}

func (s *MongoDbStorage) DeleteById(eventId primitive.ObjectID) (error) {

	filter := bson.M{"_id":eventId}

	if _, err := s.db.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return nil
}