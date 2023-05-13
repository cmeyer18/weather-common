package collections

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"go.mongodb.org/mongo-driver/bson"
)

type NotificationCollection struct {
	*custom_mongo.BaseCollection[data_structures.Notification]
}

func notificationCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "notifications"}
}

func NewNotificationCollection(connection custom_mongo.BaseConnection) (NotificationCollection, bool, bool, error) {
	collection, dbExist, collectExist, err := custom_mongo.GetCollection[data_structures.Notification](notificationCollectionArgs(), connection)
	if err != nil {
		return NotificationCollection{}, false, false, err
	}

	collectionConv := NotificationCollection{&collection}
	return collectionConv, dbExist, collectExist, nil
}

func (nc *NotificationCollection) GetFilterOnCountyOrZone(zonesList []string) ([]data_structures.Notification, error) {
	filterZone := bson.M{"zonecode": bson.M{"$in": zonesList}}
	filterCounty := bson.M{"countycode": bson.M{"$in": zonesList}}
	// Get all the records and process them into an array
	return nc.Get(bson.M{"$or": []interface{}{filterCounty, filterZone}})
}
