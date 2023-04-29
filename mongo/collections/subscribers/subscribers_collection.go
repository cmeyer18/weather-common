package subscribers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Subscriber struct {
	URL string `json:"url" bson:"url"`
}

type SubscriberCollection struct {
	subscriberCollection *mongo.Collection
	Client               *mongo.Client
	Ctx                  context.Context
}

func NewSubscriberCollection(client *mongo.Client, ctx context.Context) (SubscriberCollection, error) {
	collection := SubscriberCollection{Ctx: ctx, Client: client}

	subscriberCollection, err := getCollection(client, ctx)
	if err != nil {
		log.WithError(err).Error("mongo.subscriberCollection.getCollection.NewSubscriberCollection.Error")
		return SubscriberCollection{}, err
	}

	collection.subscriberCollection = subscriberCollection

	return collection, nil
}

// GetCollection gets a collection for the notifier service
func getCollection(client *mongo.Client, ctx context.Context) (*mongo.Collection, error) {
	// Alert if no db found
	dbNames, err := client.ListDatabaseNames(ctx, bson.M{"name": "weather"})
	if err != nil {
		log.WithError(err).Error("mongo.subscriberCollection.getCollection.ListDatabaseNames.error")
		return nil, err
	}

	if len(dbNames) == 0 {
		log.Warn("mongo.subscriberCollection.getCollection.ListDatabasesNames weather db not found")
	}

	weatherDB := client.Database("weather")

	// Alert if no collection found
	collectionNames, err := weatherDB.ListCollectionNames(ctx, bson.M{"name": "subscribers"})
	if err != nil {
		log.WithError(err).Error("mongo.subscriberCollection.getCollection.ListCollectionNames.error")
		return nil, err
	}
	if len(collectionNames) == 0 {
		log.Warn("mongo.subscriberCollection.getCollection.ListCollectionNames subscribers collection not found")
	}

	subscriberCollection := weatherDB.Collection("subscribers")
	return subscriberCollection, nil
}

func (sc *SubscriberCollection) FindAll() ([]Subscriber, error) {
	var subscribers []Subscriber

	// Get all the records and process them into an array
	results, err := sc.subscriberCollection.Find(sc.Ctx, bson.M{})
	if err != nil {
		log.WithError(err).Error("mongo.subscriberCollection.FindAll.error")
		return nil, err
	}

	for results.Next(sc.Ctx) {
		var subscriber Subscriber
		err := results.Decode(&subscriber)
		if err != nil {
			log.WithError(err).Error("mongo.subscriberCollection.FindAll.Decode.failure")
			return nil, err
		}

		subscribers = append(subscribers, subscriber)
	}

	return subscribers, nil
}

// Insert takes the given alerts and puts them in the mongo db if they aren't in there already.
func (sc *SubscriberCollection) Insert(subscribers []Subscriber) error {
	mongoInterface := make([]interface{}, len(subscribers))

	for i, subscriber := range subscribers {
		mongoInterface[i] = subscriber
	}

	_, err := sc.subscriberCollection.InsertMany(sc.Ctx, mongoInterface)
	if err != nil {
		log.WithError(err).Error("mongo.subscriberCollection.Insert.failure")
		return err
	}
	log.Info("mongo.subscriberCollection.Inserted")
	return nil
}

// Delete takes the given alerts and deletes them in the db.
func (sc *SubscriberCollection) Delete(subscribers []string) error {
	_, err := sc.subscriberCollection.DeleteMany(sc.Ctx, bson.M{"url": bson.M{"$in": subscribers}})
	if err != nil {
		log.WithError(err).Error("mongo.subscriberCollection.Delete.failure")
		return err
	}
	log.Info("mongo.subscriberCollection.Deleted")
	return nil
}
