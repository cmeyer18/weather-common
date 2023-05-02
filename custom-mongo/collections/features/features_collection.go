package features

import (
	"github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/custom-mongo/collections"
	"github.com/cmeyer18/weather-common/data_structures"
	"github.com/sirupsen/logrus"
	"time"
)

const databaseName = "weather"
const collectionName = "features"

type FeatureCollection struct {
	Collection *collections.BaseCollection[data_structures.Feature]
	logger     *logrus.Logger
}

func NewFeatureCollection(baseConnection custom_mongo.BaseConnection) (FeatureCollection, error) {
	featureCollection, err := custom_mongo.GetCollection[data_structures.Feature](databaseName, collectionName, baseConnection)
	if err != nil {
		baseConnection.Logger.WithError(err).Error("custom-mongo.collections.features_collection.NewFeatureCollection.Error")
		return FeatureCollection{}, err
	}

	return FeatureCollection{Collection: featureCollection, logger: baseConnection.Logger}, nil
}

func (fc *FeatureCollection) GetExpiredFeatures(givenTime time.Time) ([]data_structures.Feature, error) {
	var expiredAlerts []data_structures.Feature

	// Get all records, delete the ones where the time has expired
	features, err := fc.Collection.GetAll()
	if err != nil {
		fc.logger.WithError(err).Error("custom-mongo.collections.features_collection.GetAll.error")
		return nil, err
	}

	for _, feature := range features {
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
