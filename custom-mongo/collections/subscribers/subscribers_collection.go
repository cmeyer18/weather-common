package subscribers

import (
	"context"
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const databaseName = "weather"
const collectionName = "subscribers"

type SubscriberCollection struct {
	subscriberCollection *mongo.Collection
	logger               *logrus.Logger
}

func NewSubscriberCollection(baseConnection custom_mongo.BaseConnection) (SubscriberCollection, error) {
	subscriberCollection, err := baseConnection.GetCollection(databaseName, collectionName)
	if err != nil {
		baseConnection.Logger.WithError(err).Error("custom-mongo.subscribers_collection.getCollection.NewSubscriberCollection.Error")
		return SubscriberCollection{}, err
	}

	return SubscriberCollection{subscriberCollection: subscriberCollection, logger: baseConnection.Logger}, nil
}

func (sc *SubscriberCollection) FindAll() ([]data_structures.Subscriber, error) {
	var subscribers []data_structures.Subscriber

	// Get all the records and process them into an array
	results, err := sc.subscriberCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		sc.logger.WithError(err).Error("custom-mongo.subscribers_collection.FindAll.error")
		return nil, err
	}

	for results.Next(context.TODO()) {
		var subscriber data_structures.Subscriber
		err := results.Decode(&subscriber)
		if err != nil {
			sc.logger.WithError(err).Error("custom-mongo.subscribers_collection.FindAll.Decode.failure")
			return nil, err
		}

		subscribers = append(subscribers, subscriber)
	}

	return subscribers, nil
}

// Insert takes the given alerts and puts them in the custom-mongo db if they aren't in there already.
func (sc *SubscriberCollection) Insert(subscribers []data_structures.Subscriber) error {
	mongoInterface := make([]interface{}, len(subscribers))

	for i, subscriber := range subscribers {
		mongoInterface[i] = subscriber
	}

	_, err := sc.subscriberCollection.InsertMany(context.TODO(), mongoInterface)
	if err != nil {
		sc.logger.WithError(err).Error("custom-mongo.subscribers_collection.Insert.failure")
		return err
	}
	sc.logger.Info("custom-mongo.subscribers_collection.Inserted")
	return nil
}

// Delete takes the given alerts and deletes them in the db.
func (sc *SubscriberCollection) Delete(subscribers []string) error {
	_, err := sc.subscriberCollection.DeleteMany(context.TODO(), bson.M{"url": bson.M{"$in": subscribers}})
	if err != nil {
		sc.logger.WithError(err).Error("custom-mongo.subscribers_collection.Delete.failure")
		return err
	}
	sc.logger.Info("custom-mongo.subscribers_collection.Deleted")
	return nil
}
