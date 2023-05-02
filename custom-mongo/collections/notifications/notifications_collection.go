package subscribers

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/custom-mongo/collections"
	"github.com/cmeyer18/weather-common/data_structures"
	"github.com/sirupsen/logrus"
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
