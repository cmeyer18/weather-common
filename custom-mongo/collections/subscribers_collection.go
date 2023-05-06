package collections

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
)

type subscriberCollectionI custom_mongo.BaseCollection[data_structures.Subscriber]

type SubscriberCollection struct {
	subscriberCollectionI
}

func subscriberCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "subscribers"}
}

func NewSubscriberCollection(connection custom_mongo.BaseConnection) (SubscriberCollection, bool, bool, error) {
	collection, dbExist, collectExist, err := custom_mongo.GetCollection[data_structures.Subscriber](subscriberCollectionArgs(), connection)
	if err != nil {
		return SubscriberCollection{}, false, false, err
	}

	collectionConv := SubscriberCollection{subscriberCollectionI(collection)}

	return collectionConv, dbExist, collectExist, nil
}
