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

func (bc *BaseCollection[T]) Exists(id string) (bool, error) {
	result := bc.collection.FindOne(context.TODO(), bson.M{"id": id})

	var element T
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

func (bc *BaseCollection[T]) Delete(ids []string) error {
	_, err := bc.collection.DeleteMany(context.TODO(), bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return err
	}
	return nil
}

func (bc *BaseCollection[T]) Get(id string) (T, error) {
	var element T

	// Get all the records and process them into an array
	result := bc.collection.FindOne(context.TODO(), bson.M{"id": id})

	err := result.Decode(&element)
	if err == mongo.ErrNoDocuments {
		return element, nil
	} else if err != nil {
		return element, err
	}

	return element, nil
}

func (bc *BaseCollection[T]) GetManyWithFilter(bsonFilter bson.M) ([]T, error) {
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

// GetCollection gets a collection for the notifier service
func GetCollection[T any](args BaseCollectionArgs, connection BaseConnection) (BaseCollection[T], bool, bool, error) {
	databaseFound := true
	collectionFound := true

	// Alert if no db found
	dbFoundList, err := connection.client.ListDatabaseNames(context.TODO(), bson.M{"name": args.DatabaseName})
	if err != nil {
		return BaseCollection[T]{}, false, false, err
	}
	if len(dbFoundList) == 0 {
		databaseFound = false
	}

	database := connection.client.Database(args.DatabaseName)

	// Alert if no collection found
	collectionNames, err := database.ListCollectionNames(context.TODO(), bson.M{"name": args.CollectionName})
	if err != nil {
		return BaseCollection[T]{}, false, false, err
	}
	if len(collectionNames) == 0 {
		collectionFound = true
	}

	collection := database.Collection(args.CollectionName)
	baseCollection := NewBaseCollection[T](collection)

	return baseCollection, databaseFound, collectionFound, nil
}
