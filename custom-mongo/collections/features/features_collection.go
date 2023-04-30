package features

import (
	"context"
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const databaseName = "weather"
const collectionName = "features"

type FeatureCollection struct {
	featureCollection *mongo.Collection
	logger            *logrus.Logger
}

func NewFeatureCollection(baseConnection custom_mongo.BaseConnection) (FeatureCollection, error) {
	featureCollection, err := baseConnection.GetCollection(databaseName, collectionName)
	if err != nil {
		baseConnection.Logger.WithError(err).Error("custom-mongo.collections.features_collection.NewFeatureCollection.Error")
		return FeatureCollection{}, err
	}

	return FeatureCollection{featureCollection: featureCollection, logger: baseConnection.Logger}, nil
}

func (fc *FeatureCollection) Exists(id string) (bool, error) {
	result := fc.featureCollection.FindOne(context.TODO(), bson.M{"id": id})

	var feature data_structures.Feature
	err := result.Decode(&feature)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.Exists.failure error with decoding")
		return false, err
	}

	return true, nil
}

func (fc *FeatureCollection) GetExpiredFeatures(givenTime time.Time) ([]data_structures.Feature, error) {
	var expiredAlerts []data_structures.Feature

	// Get all records, delete the ones where the time has expired
	results, err := fc.featureCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.GetExpiredFeatures.Find.error")
		return nil, err
	}

	for results.Next(context.TODO()) {

		// Convert to a feature
		var feature data_structures.Feature
		err := results.Decode(&feature)
		if err != nil {
			fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.GetExpiredFeatures.Decode.failure")
			return nil, err
		}

		var definedTime time.Time

		// If there is an end time, start by using this time.
		if feature.Properties.Ends != "" {
			endingTime, err := time.Parse(time.RFC3339, feature.Properties.Ends)
			if err != nil {
				fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.GetExpiredFeatures.Parse.failed parsing ending time")
				return nil, err
			}
			definedTime = endingTime
		}

		// If there is an expired time, lets process it as well.
		if feature.Properties.Expires != "" {
			expiredTime, err := time.Parse(time.RFC3339, feature.Properties.Expires)
			if err != nil {
				fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.GetExpiredFeatures.Parse.failed parsing expired time")
				return nil, err
			}
			if definedTime.IsZero() {
				definedTime = expiredTime
			} else {
				if expiredTime.After(definedTime) {
					definedTime = expiredTime
				}
			}

		}

		// Delete the alert if it is before now
		if definedTime.Before(givenTime) {
			fc.logger.Info("custom-mongo.collections.features_collection.GetExpiredFeatures alert expired at " + definedTime.Format(time.UnixDate) + " , removing")
			expiredAlerts = append(expiredAlerts, feature)
		}
	}

	return expiredAlerts, nil
}

// Insert takes the given alerts and puts them in the custom-mongo db if they aren't in there already.
func (fc *FeatureCollection) Insert(features []data_structures.Feature) error {
	mongoInterface := make([]interface{}, len(features))

	for i, feature := range features {
		mongoInterface[i] = feature
	}

	_, err := fc.featureCollection.InsertMany(context.TODO(), mongoInterface)
	if err != nil {
		fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.Insert.Error")
		return err
	}
	fc.logger.Info("custom-mongo.collections.features_collection.Insert.Inserted")
	return nil
}

// Delete takes the given alerts and deletes them in the db.
func (fc *FeatureCollection) Delete(featureIds []string) error {
	_, err := fc.featureCollection.DeleteMany(context.TODO(), bson.M{"id": bson.M{"$in": featureIds}})
	if err != nil {
		fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.Delete.failure")
		return err
	}
	fc.logger.Info("custom-mongo.collections.features_collection.Delete.Deleted")
	return nil
}
