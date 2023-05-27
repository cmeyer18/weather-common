package collections

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
)

type ServerNotificationSubscriberCollection struct {
	*custom_mongo.BaseCollection[data_structures.ServerNotificationSubscriber]
}

func serverNotificationSubscriberCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "subscribers"}
}

func NewServerNotificationSubscriberCollection(connection custom_mongo.BaseConnection) (ServerNotificationSubscriberCollection, bool, bool, error) {
	collection, dbExist, collectExist, err := custom_mongo.GetCollection[data_structures.ServerNotificationSubscriber](serverNotificationSubscriberCollectionArgs(), connection)
	if err != nil {
		return ServerNotificationSubscriberCollection{}, false, false, err
	}

	collectionConv := ServerNotificationSubscriberCollection{&collection}
	return collectionConv, dbExist, collectExist, nil
}
