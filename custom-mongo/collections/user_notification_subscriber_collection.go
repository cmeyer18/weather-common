package collections

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"go.mongodb.org/mongo-driver/bson"
)

type UserNotificationSubscriberCollection struct {
	*custom_mongo.BaseCollection[data_structures.UserNotificationSubscriber]
}

func userNotificationSubscriberCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "notifications"}
}

func NewUserNotificationSubscriberCollection(connection custom_mongo.BaseConnection) (UserNotificationSubscriberCollection, bool, bool, error) {
	collection, dbExist, collectExist, err := custom_mongo.GetCollection[data_structures.UserNotificationSubscriber](userNotificationSubscriberCollectionArgs(), connection)
	if err != nil {
		return UserNotificationSubscriberCollection{}, false, false, err
	}

	collectionConv := UserNotificationSubscriberCollection{&collection}
	return collectionConv, dbExist, collectExist, nil
}

func (nc *UserNotificationSubscriberCollection) GetFilterOnCountyOrZone(zonesList []string) ([]data_structures.UserNotificationSubscriber, error) {
	filterZone := bson.M{"zonecode": bson.M{"$in": zonesList}}
	filterCounty := bson.M{"countycode": bson.M{"$in": zonesList}}
	// Get all the records and process them into an array
	return nc.Get(bson.M{"$or": []interface{}{filterCounty, filterZone}})
}
