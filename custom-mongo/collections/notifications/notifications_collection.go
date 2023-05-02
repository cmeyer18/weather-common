package subscribers

import (
	"context"
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/custom-mongo/collections"
	"github.com/cmeyer18/weather-common/data_structures"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const databaseName = "weather"
const collectionName = "notifications"

type NotificationCollection struct {
	Collection *collections.BaseCollection[data_structures.Notification]
	logger     *logrus.Logger
}

func NewNotificationCollection(baseConnection custom_mongo.BaseConnection) (NotificationCollection, error) {
	notificationCollection, err := custom_mongo.GetCollection[data_structures.Notification](databaseName, collectionName, baseConnection)
	if err != nil {
		baseConnection.Logger.WithError(err).Error("custom-mongo.collections.notification_collection.NewNotificationCollection.Error")
		return NotificationCollection{}, err
	}

	return NotificationCollection{Collection: notificationCollection, logger: baseConnection.Logger}, nil
}

func (nc *NotificationCollection) GetFilterOnCountyOrZone(zonesList []string) ([]data_structures.Notification, error) {
	var notifications []data_structures.Notification

	filterZone := bson.M{"zonecode": bson.M{"$in": zonesList}}
	filterCounty := bson.M{"countycode": bson.M{"$in": zonesList}}
	// Get all the records and process them into an array
	results, err := nc.Collection.Collection.Find(context.TODO(), bson.M{"$or": []interface{}{filterCounty, filterZone}})
	if err != nil {
		nc.logger.WithError(err).Error("custom-mongo.notification_collection.GetFilterOnCountyOrZone.Find.error")
		return nil, err
	}

	for results.Next(context.TODO()) {
		var element data_structures.Notification
		err := results.Decode(&element)
		if err != nil {
			nc.logger.WithError(err).Error("custom-mongo.notification_collection.GetFilterOnCountyOrZone.Decode.failure")
			return nil, err
		}

		notifications = append(notifications, element)
	}

	return notifications, nil
}
