package custom_mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseCollection[T any] struct {
	collection *mongo.Collection
}

func NewBaseCollection[T any](collection *mongo.Collection) BaseCollection[T] {
	return BaseCollection[T]{collection: collection}
}

func (bc *BaseCollection[T]) Exists(bsonFilter bson.M) (bool, error) {
	var element T

	result := bc.collection.FindOne(context.TODO(), bsonFilter)
	err := result.Decode(&element)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (bc *BaseCollection[T]) Insert(elements []T) error {
	mongoInterface := make([]interface{}, len(elements))

	for i, element := range elements {
		mongoInterface[i] = element
	}

	_, err := bc.collection.InsertMany(context.TODO(), mongoInterface)
	if err != nil {
		return err
	}

	return nil
}

func (bc *BaseCollection[T]) Delete(bsonFilter bson.M) error {
	_, err := bc.collection.DeleteMany(context.TODO(), bsonFilter)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BaseCollection[T]) Get(bsonFilter bson.M) ([]T, error) {
	var elements []T

	// Get all the records and process them into an array
	results, err := bc.collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}

	for results.Next(context.TODO()) {
		var element T
		err = results.Decode(&element)
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)
	}

	return elements, nil
}

func (bc *BaseCollection[T]) GetAll() ([]T, error) {
	var elements []T

	// Get all the records and process them into an array
	results, err := bc.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	for results.Next(context.TODO()) {
		var element T
		err := results.Decode(&element)
		if err != nil {
			return nil, err
		}

		elements = append(elements, element)
	}

	return elements, nil
}
