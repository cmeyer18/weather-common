package subscribers

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/custom-mongo/collections"
	"github.com/cmeyer18/weather-common/data_structures"
	"github.com/sirupsen/logrus"
)

const databaseName = "weather"
const collectionName = "subscribers"

type SubscriberCollection struct {
	Collection *collections.BaseCollection[data_structures.Subscriber]
	logger     *logrus.Logger
}

func NewSubscriberCollection(baseConnection custom_mongo.BaseConnection) (SubscriberCollection, error) {
	collection, err := custom_mongo.GetCollection[data_structures.Subscriber](databaseName, collectionName, baseConnection)
	if err != nil {
		baseConnection.Logger.WithError(err).Error("custom-mongo.collections.subscriber_collection.NewSubscriberCollection.Error")
		return SubscriberCollection{}, err
	}

	return SubscriberCollection{Collection: collection, logger: baseConnection.Logger}, nil
}
