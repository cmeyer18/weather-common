package collections

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
)

type SubscriberCollection struct {
	custom_mongo.BaseCollection[data_structures.Subscriber]
}

func NewSubscriberCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "subscribers"}
}
