package collections

import (
	custommongo "github.com/cmeyer18/weather-common/custom-mongo"
	"github.com/cmeyer18/weather-common/data_structures"
	"time"
)

type AlertCollection struct {
	*custommongo.BaseCollection[data_structures.Alert]
}

func NewAlertCollection(connection custommongo.BaseConnection) (AlertCollection, bool, bool, error) {
	collection, dbExist, collectExist, err := custommongo.GetCollection[data_structures.Alert](alertCollectionArgs(), connection)
	if err != nil {
		return AlertCollection{}, false, false, err
	}

	collectionConv := AlertCollection{&collection}

	return collectionConv, dbExist, collectExist, nil
}

func alertCollectionArgs() custommongo.BaseCollectionArgs {
	return custommongo.BaseCollectionArgs{DatabaseName: "weather", CollectionName: "alerts"}
}

func (ac *AlertCollection) GetExpiredAlerts(givenTime time.Time) ([]data_structures.Alert, error) {
	var expiredAlerts []data_structures.Alert

	// Get all records, delete the ones where the time has expired
	alerts, err := ac.GetAll()
	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		var definedTime time.Time

		if alert.Properties.Expires != "" {
			expiredTime, err := time.Parse(time.RFC3339, alert.Properties.Expires)
			if err != nil {
				return nil, err
			}
			definedTime = expiredTime
		}

		if alert.Properties.Ends != "" {
			endingTime, err := time.Parse(time.RFC3339, alert.Properties.Ends)
			if err != nil {
				return nil, err
			}
			if definedTime.IsZero() {
				definedTime = endingTime
			} else {
				if definedTime.Before(endingTime) {
					definedTime = endingTime
				}
			}
		}

		// Delete the alert if it is before now
		if definedTime.Before(givenTime) {
			expiredAlerts = append(expiredAlerts, alert)
		}
	}

	return expiredAlerts, nil
}
