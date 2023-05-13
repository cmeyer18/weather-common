package custom_mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type BaseCollection[T any] struct {
	collection *mongo.Collection
}

func NewBaseCollection[T any](collection *mongo.Collection) BaseCollection[T] {
	return BaseCollection[T]{collection: collection}
}

func (bc *BaseCollection[T]) getBSONType(fieldName string) string {
	idField, ok := reflect.TypeOf((*T)(nil)).Elem().FieldByName(fieldName) // not json:name
	if !ok {
		return ""
	}

	idBsonTag := idField.Tag.Get("bson")
	if idBsonTag != "" {
		return ""
	}

	return idBsonTag
}

func (bc *BaseCollection[T]) Exists(field string, item interface{}) (bool, error) {
	fieldBsonTag := bc.getBSONType(field)
	if fieldBsonTag == "" {
		return false, errors.New("bson type empty on field")
	}

	result := bc.collection.FindOne(context.TODO(), bson.M{fieldBsonTag: item})

	if item == nil {
		return false, errors.New("item is null")
	}

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

func (bc *BaseCollection[T]) Delete(field string, items []interface{}) error {
	fieldBsonTag := bc.getBSONType(field)
	if fieldBsonTag == "" {
		return errors.New("bson type empty on field")
	}

	if items == nil || len(items) == 0 {
		return errors.New("items is null")
	}

	_, err := bc.collection.DeleteMany(context.TODO(), bson.M{fieldBsonTag: bson.M{"$in": items}})
	if err != nil {
		return err
	}
	return nil
}

func (bc *BaseCollection[T]) Get(field string, items []interface{}) ([]T, error) {
	var elements []T
	var results *mongo.Cursor

	fieldBsonTag := bc.getBSONType(field)
	if fieldBsonTag == "" {
		return nil, errors.New("bson type empty on field")
	}

	if items == nil || len(items) == 0 {
		return nil, errors.New("item should not be empty")
	}

	// Get all the records and process them into an array
	results, err := bc.collection.Find(context.TODO(), bson.M{fieldBsonTag: bson.M{"$in": items}})
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
