package repo

import (
	"client_service/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Config *config.MongoConfig
	Client *mongo.Client
}

func (MongoDB *MongoDB) AddTrip(trip Trip) error {
	collection := MongoDB.Client.Database("client_service").Collection("trips")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, trip)
	if err != nil {
		return err
	}
	return nil
}

func (MongoDB *MongoDB) GetTrips(tripId string) ([]Trip, error) {
	collection := MongoDB.Client.Database("client_service").Collection("trips")

	filter := bson.M{"trip_id": tripId}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return []Trip{}, fmt.Errorf(err.Error())
	}

	var result []Trip
	if err := cursor.All(context.Background(), &result); err != nil {
		return []Trip{}, fmt.Errorf(err.Error())
	}
	return result, nil
}

func (MongoDB *MongoDB) GetTripByTripId(tripId string) (Trip, error) {
	collection := MongoDB.Client.Database("client_service").Collection("trips")

	filter := bson.M{"trip_id": tripId}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return Trip{}, fmt.Errorf(err.Error())
	}

	var result []Trip
	if err := cursor.All(context.Background(), &result); err != nil {
		log.Fatal(err)
		return Trip{}, fmt.Errorf("MongaDB Error")
	}
	if len(result) == 0 {
		return Trip{}, fmt.Errorf("Not Found")
	}
	return result[0], nil
}

func (MongoDB *MongoDB) GetTripByOfferId(offerId string) (Trip, error) {
	collection := MongoDB.Client.Database("client_service").Collection("trips")

	filter := bson.M{"offer_id": offerId}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return Trip{}, fmt.Errorf(err.Error())
	}

	var result []Trip
	if err := cursor.All(context.Background(), &result); err != nil {
		log.Fatal(err)
		return Trip{}, fmt.Errorf("MongaDB Error")
	}
	if len(result) == 0 {
		return Trip{}, fmt.Errorf("Not Found")
	}
	return result[0], nil
}

func (MongoDB *MongoDB) ChangeTripByTripId(tripId string, key string, value interface{}) error {
	collection := MongoDB.Client.Database("client_service").Collection("trips")

	filter := bson.M{"trip_id": tripId}

	update := bson.M{"$set": bson.M{key: value}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (MongoDB *MongoDB) ChangeTripByOfferId(offerId string, key string, value interface{}) error {
	collection := MongoDB.Client.Database("client_service").Collection("trips")

	filter := bson.M{"offer_id": offerId}

	update := bson.M{"$set": bson.M{key: value}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (MongoDB *MongoDB) Serve() error {
	clientOptions := options.Client().ApplyURI(MongoDB.Config.URI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Connected to MongoDB!")
	MongoDB.Client = client

	database := client.Database("client_service")
	collectionName := "trips"
	collectionOptions := options.CreateCollection()
	err = database.CreateCollection(context.Background(), collectionName, collectionOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}
