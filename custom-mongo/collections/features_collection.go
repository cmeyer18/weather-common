package collections

import (
	custom_mongo "github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"time"
)

type FeatureCollection struct {
	custom_mongo.BaseCollection[data_structures.Feature]
}

func NewFeatureCollectionArgs() custom_mongo.BaseCollectionArgs {
	return custom_mongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "features"}
}

func (fc *FeatureCollection) GetExpiredFeatures(givenTime time.Time) ([]data_structures.Feature, error) {
	var expiredAlerts []data_structures.Feature

	// Get all records, delete the ones where the time has expired
	features, err := fc.GetAll()
	if err != nil {
		return nil, err
	}

	for _, feature := range features {
		var definedTime time.Time

		// If there is an end time, start by using this time.
		if feature.Properties.Ends != "" {
			endingTime, err := time.Parse(time.RFC3339, feature.Properties.Ends)
			if err != nil {
				return nil, err
			}
			definedTime = endingTime
		}

		// If there is an expired time, lets process it as well.
		if feature.Properties.Expires != "" {
			expiredTime, err := time.Parse(time.RFC3339, feature.Properties.Expires)
			if err != nil {
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
			expiredAlerts = append(expiredAlerts, feature)
		}
	}

	return expiredAlerts, nil
}
