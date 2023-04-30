package subscribers

import (
	"context"
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const databaseName = "weather"
const collectionName = "notifications"

type NotificationCollection struct {
	notificationCollection *mongo.Collection
	logger                 *logrus.Logger
}

func NewNotificationCollection(baseConnection custom_mongo.BaseConnection) (NotificationCollection, error) {
	notificationCollection, err := baseConnection.GetCollection(databaseName, collectionName)
	if err != nil {
		baseConnection.Logger.WithError(err).Error("custom-mongo.notifications_collection.NewNotificationCollection.Error")
		return NotificationCollection{}, err
	}

	return NotificationCollection{notificationCollection: notificationCollection, logger: baseConnection.Logger}, nil
}

func (nc *NotificationCollection) FindAll() ([]data_structures.Notification, error) {
	var subscribers []data_structures.Notification

	// Get all the records and process them into an array
	results, err := nc.notificationCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		nc.logger.WithError(err).Error("custom-mongo.notifications_collection.FindAll.error")
		return nil, err
	}

	for results.Next(context.TODO()) {
		var notification data_structures.Notification
		err := results.Decode(&notification)
		if err != nil {
			nc.logger.WithError(err).Error("custom-mongo.notifications_collection.FindAll.Decode.failure")
			return nil, err
		}

		subscribers = append(subscribers, notification)
	}

	return subscribers, nil
}

// Insert takes the given alerts and puts them in the custom-mongo db if they aren't in there already.
func (nc *NotificationCollection) Insert(subscribers []NotificationCollection) error {
	mongoInterface := make([]interface{}, len(subscribers))

	for i, subscriber := range subscribers {
		mongoInterface[i] = subscriber
	}

	_, err := nc.notificationCollection.InsertMany(context.TODO(), mongoInterface)
	if err != nil {
		nc.logger.WithError(err).Error("custom-mongo.notifications_collection.Insert.failure")
		return err
	}
	nc.logger.Info("custom-mongo.notifications_collection.Inserted")
	return nil
}

// Delete takes the given alerts and deletes them in the db.
func (nc *NotificationCollection) Delete(notificationIds []string) error {
	_, err := nc.notificationCollection.DeleteMany(context.TODO(), bson.M{"notificationid": bson.M{"$in": notificationIds}})
	if err != nil {
		nc.logger.WithError(err).Error("custom-mongo.notifications_collection.Delete.failure")
		return err
	}
	nc.logger.Info("custom-mongo.notifications_collection.Deleted")
	return nil
}
