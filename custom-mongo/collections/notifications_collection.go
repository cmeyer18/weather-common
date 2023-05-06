package collections

import (
	"context"
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"go.mongodb.org/mongo-driver/bson"
)

type NotificationCollection struct {
	custom_mongo.BaseCollection[data_structures.Notification]
}

func NewNotificationCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "notifications"}
}

func (nc *NotificationCollection) GetFilterOnCountyOrZone(zonesList []string) ([]data_structures.Notification, error) {
	var notifications []data_structures.Notification

	filterZone := bson.M{"zonecode": bson.M{"$in": zonesList}}
	filterCounty := bson.M{"countycode": bson.M{"$in": zonesList}}
	// Get all the records and process them into an array
	results, err := nc.Collection.Find(context.TODO(), bson.M{"$or": []interface{}{filterCounty, filterZone}})
	if err != nil {
		return nil, err
	}

	for results.Next(context.TODO()) {
		var element data_structures.Notification
		err := results.Decode(&element)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, element)
	}

	return notifications, nil
}
