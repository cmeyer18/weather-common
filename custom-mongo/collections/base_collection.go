package collections

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseCollectionInterface[T any] interface {
	Exists(id string) (bool, error)
	Insert(elements []T) error
	Delete(ids []string) error
	Get(id string) (T, error)
	GetAll() ([]T, error)
}

type BaseCollection[T any] struct {
	Collection *mongo.Collection
	logger     *logrus.Logger
}

func NewBaseCollection[T any](collection *mongo.Collection, logger *logrus.Logger) BaseCollection[T] {
	return BaseCollection[T]{Collection: collection, logger: logger}
}

func (bc *BaseCollection[T]) Exists(id string) (bool, error) {
	result := bc.Collection.FindOne(context.TODO(), bson.M{"id": id})

	var element T
	err := result.Decode(&element)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		bc.logger.WithError(err).Error("custom-mongo.collections.base_collection.Exists.failure error with decoding")
		return false, err
	}

	return true, nil
}

func (bc *BaseCollection[T]) Insert(elements []T) error {
	mongoInterface := make([]interface{}, len(elements))

	for i, element := range elements {
		mongoInterface[i] = element
	}

	_, err := bc.Collection.InsertMany(context.TODO(), mongoInterface)
	if err != nil {
		bc.logger.WithError(err).Error("custom-mongo.collections.base_collection.Insert.Error")
		return err
	}
	bc.logger.Info("custom-mongo.collections.base_collection.Insert.Inserted")
	return nil
}

func (bc *BaseCollection[T]) Delete(ids []string) error {
	_, err := bc.Collection.DeleteMany(context.TODO(), bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		bc.logger.WithError(err).Error("custom-mongo.collections.features_collection.Delete.failure")
		return err
	}
	bc.logger.Info("custom-mongo.collections.features_collection.Delete.Deleted")
	return nil
}

func (bc *BaseCollection[T]) Get(id string) (T, error) {
	var element T

	// Get all the records and process them into an array
	result := bc.Collection.FindOne(context.TODO(), bson.M{"id": id})

	err := result.Decode(&element)
	if err == mongo.ErrNoDocuments {
		return element, nil
	} else if err != nil {
		bc.logger.WithError(err).Error("custom-mongo.collections.base_collection.Exists.failure error with decoding")
		return element, err
	}

	return element, nil
}

func (bc *BaseCollection[T]) GetAll() ([]T, error) {
	var elements []T

	// Get all the records and process them into an array
	results, err := bc.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		bc.logger.WithError(err).Error("custom-mongo.subscribers_collection.FindAll.error")
		return nil, err
	}

	for results.Next(context.TODO()) {
		var element T
		err := results.Decode(&element)
		if err != nil {
			bc.logger.WithError(err).Error("custom-mongo.subscribers_collection.FindAll.Decode.failure")
			return nil, err
		}

		elements = append(elements, element)
	}

	return elements, nil
}
