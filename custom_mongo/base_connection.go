package custom_mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type BaseCollectionArgs struct {
	DatabaseName   string
	CollectionName string
}

// BaseConnection details a connection struct for the main file.
type BaseConnection struct {
	credentials        options.Credential
	mongoConnectionURL string
	client             *mongo.Client
}

// NewConnection returns a BaseConnection for all mongodb transactions
func NewConnection(credentials options.Credential, url string) BaseConnection {
	return BaseConnection{credentials: credentials, mongoConnectionURL: url}
}

// Ping tries to connect the server.
func (bc *BaseConnection) ping() error {
	if err := bc.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return err
	}
	return nil
}

// Disconnect takes the given ctx and client and disconnects them.
func (bc *BaseConnection) Disconnect() error {
	return bc.client.Disconnect(context.TODO())
}

// Connect sets up a connection to a mongoDB server.
func (bc *BaseConnection) Connect() error {
	connectionOptions := options.Client()
	connectionOptions.ApplyURI("mongodb://" + bc.mongoConnectionURL)
	connectionOptions.SetAuth(bc.credentials)

	client, err := mongo.Connect(context.TODO(), connectionOptions)
	if err != nil {
		return err
	}
	bc.client = client

	// Ping for up status
	err = bc.ping()
	if err != nil {
		err2 := bc.Disconnect()
		if err2 != nil {
			return err2
		}
		return err
	}

	return nil
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
