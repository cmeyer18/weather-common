package features

import (
	"context"
	"github.com/cmeyer18/weather-common/data_structures"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type FeatureCollection struct {
	featureCollection *mongo.Collection
	Client            *mongo.Client
	Ctx               context.Context
}

func NewFeatureCollection(client *mongo.Client, ctx context.Context) (FeatureCollection, error) {
	collection := FeatureCollection{Ctx: ctx, Client: client}

	featureCollection, err := getCollection(client, ctx)
	if err != nil {
		log.WithError(err).Error("mongo.collections.features_collection.NewFeatureCollection.Error")
		return FeatureCollection{}, err
	}

	collection.featureCollection = featureCollection

	return collection, nil
}

// getCollection gets a collection for the notifier service
func getCollection(client *mongo.Client, ctx context.Context) (*mongo.Collection, error) {
	// Alert if no db found
	dbNames, err := client.ListDatabaseNames(ctx, bson.M{"name": "weather"})
	if err != nil {
		log.WithError(err).Error("mongo.collections.features_collection.getCollection.ListDatabaseNames.error")
		return nil, err
	}

	if len(dbNames) == 0 {
		log.Warn("mongo.collections.features_collection.getCollection.ListDatabasesNames weather db not found")
	}

	weatherDB := client.Database("weather")

	// Alert if no collection found
	collectionNames, err := weatherDB.ListCollectionNames(ctx, bson.M{"name": "features"})
	if err != nil {
		log.WithError(err).Error("mongo.collections.features_collection.getCollection.ListCollectionNames.error")
		return nil, err
	}
	if len(collectionNames) == 0 {
		log.Warn("mongo.collections.features_collection.getCollection.ListCollectionNames features collection not found")
	}

	featuresCollection := weatherDB.Collection("features")
	return featuresCollection, nil
}

func (fc *FeatureCollection) Exists(id string) (bool, error) {
	result := fc.featureCollection.FindOne(fc.Ctx, bson.M{"id": id})

	var feature data_structures.Feature
	err := result.Decode(&feature)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		log.WithError(err).Error("mongo.collections.features_collection.Exists.failure error with decoding")
		return false, err
	}

	return true, nil
}

func (fc *FeatureCollection) GetExpiredFeatures(givenTime time.Time) ([]data_structures.Feature, error) {
	var expiredAlerts []data_structures.Feature

	// Get all records, delete the ones where the time has expired
	results, err := fc.featureCollection.Find(fc.Ctx, bson.M{})
	if err != nil {
		log.WithError(err).Error("mongo.collections.features_collection.GetExpiredFeatures.Find.error")
		return nil, err
	}

	for results.Next(fc.Ctx) {

		// Convert to a feature
		var feature data_structures.Feature
		err := results.Decode(&feature)
		if err != nil {
			log.WithError(err).Error("mongo.collections.features_collection.GetExpiredFeatures.Decode.failure")
			return nil, err
		}

		var definedTime time.Time

		// If there is an end time, start by using this time.
		if feature.Properties.Ends != "" {
			endingTime, err := time.Parse(time.RFC3339, feature.Properties.Ends)
			if err != nil {
				log.WithError(err).Error("mongo.collections.features_collection.GetExpiredFeatures.Parse.failed parsing ending time")
				return nil, err
			}
			definedTime = endingTime
		}

		// If there is an expired time, lets process it as well.
		if feature.Properties.Expires != "" {
			expiredTime, err := time.Parse(time.RFC3339, feature.Properties.Expires)
			if err != nil {
				log.WithError(err).Error("mongo.collections.features_collection.GetExpiredFeatures.Parse.failed parsing expired time")
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
			log.Info("mongo.collections.features_collection.GetExpiredFeatures alert expired at " + definedTime.Format(time.UnixDate) + " , removing")
			expiredAlerts = append(expiredAlerts, feature)
		}
	}

	return expiredAlerts, nil
}

// Insert takes the given alerts and puts them in the mongo db if they aren't in there already.
func (fc *FeatureCollection) Insert(features []data_structures.Feature) error {
	mongoInterface := make([]interface{}, len(features))

	for i, feature := range features {
		mongoInterface[i] = feature
	}

	_, err := fc.featureCollection.InsertMany(fc.Ctx, mongoInterface)
	if err != nil {
		log.WithError(err).Error("mongo.collections.features_collection.Insert.Error")
		return err
	}
	log.Info("mongo.collections.features_collection.Insert.Inserted")
	return nil
}

// Delete takes the given alerts and deletes them in the db.
func (fc *FeatureCollection) Delete(featureIds []string) error {
	_, err := fc.featureCollection.DeleteMany(fc.Ctx, bson.M{"id": bson.M{"$in": featureIds}})
	if err != nil {
		log.WithError(err).Error("mongo.collections.features_collection.Delete.failure")
		return err
	}
	log.Info("mongo.collections.features_collection.Delete.Deleted")
	return nil
}
