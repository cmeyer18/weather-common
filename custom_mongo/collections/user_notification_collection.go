package collections

import (
	"github.com/cmeyer18/weather-common/custom_mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"go.mongodb.org/mongo-driver/bson"
)

type UserNotificationCollection struct {
	*custom_mongo.BaseCollection[data_structures.UserNotification]
}

func userNotificationCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "usernotification"}
}

func NewUserNotificationCollection(connection custom_mongo.BaseConnection) (UserNotificationCollection, bool, bool, error) {
	collection, dbExist, collectExist, err := custom_mongo.GetCollection[data_structures.UserNotification](userNotificationCollectionArgs(), connection)
	if err != nil {
		return UserNotificationCollection{}, false, false, err
	}

	collectionConv := UserNotificationCollection{&collection}
	return collectionConv, dbExist, collectExist, nil
}

func (nc *UserNotificationCollection) GetFilterOnCountyOrZone(zonesList []string) ([]data_structures.UserNotification, error) {
	filterZone := bson.M{"zonecode": bson.M{"$in": zonesList}}
	filterCounty := bson.M{"countycode": bson.M{"$in": zonesList}}
	// Get all the records and process them into an array
	return nc.Get(bson.M{"$or": []interface{}{filterCounty, filterZone}})
}
