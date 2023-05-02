package custom_mongo

import (
	"context"
	"github.com/cmeyer18/weather-common/custom-mongo/collections"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// BaseConnection details a connection struct for the main file.
type BaseConnection struct {
	Logger             *logrus.Logger
	credentials        options.Credential
	mongoConnectionURL string
	client             *mongo.Client
}

func NewConnection(logger *logrus.Logger, credentials options.Credential, url string, client *mongo.Client) BaseConnection {
	return BaseConnection{Logger: logger, credentials: credentials, mongoConnectionURL: url, client: client}
}

// Ping tries to connect the server.
func (bc *BaseConnection) ping() bool {
	if err := bc.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		bc.Logger.WithError(err).Error("custom-mongo.base_connection.ping.error")
		return false
	}
	bc.Logger.Info("custom-mongo.base_connection.ping.success")
	return true
}

// Disconnect takes the given ctx and client and disconnects them.
func (bc *BaseConnection) Disconnect() {
	err := bc.client.Disconnect(context.TODO())
	if err != nil {
		bc.Logger.WithError(err).Error("custom-mongo.base_connection.disconnectClient.error")
	}
	bc.Logger.Info("custom-mongo.base_connection.disconnectClient.success")
}

// Connect sets up a connection to a mongoDB server.
func (bc *BaseConnection) Connect() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+bc.mongoConnectionURL).SetAuth(bc.credentials))
	if err != nil {
		bc.Logger.WithError(err).Error("custom-mongo.base_connection.connectedClient.Connect.error")
		return err
	}
	bc.client = client

	// Ping for up status
	success := bc.ping()
	if !success {
		bc.Logger.Error("custom-mongo.base_connection.connectedClient.Ping.error")
		bc.Disconnect()
		return nil
	}

	bc.Logger.Info("custom-mongo.base_connection.connectedClient.Connect.success")
	return nil
}

// GetCollection gets a collection for the notifier service
func GetCollection[T any](databaseName, collectionName string, connection BaseConnection) (*collections.BaseCollection[T], error) {
	fields := logrus.Fields{"database": databaseName, "collection": collectionName}

	// Alert if no db found
	dbFoundList, err := connection.client.ListDatabaseNames(context.TODO(), bson.M{"name": databaseName})
	if err != nil {
		connection.Logger.WithError(err).Error("custom-mongo.base_connection.getCollection.ListDatabaseNames.error")
		return nil, err
	}
	if len(dbFoundList) == 0 {
		connection.Logger.WithFields(fields).Warn("custom-mongo.base_connection.getCollection.ListDatabasesNames db not found")
	}

	database := connection.client.Database(databaseName)

	// Alert if no collection found
	collectionNames, err := database.ListCollectionNames(context.TODO(), bson.M{"name": collectionName})
	if err != nil {
		connection.Logger.WithError(err).Error("custom-mongo.base_connection.getCollection.ListCollectionNames.error")
		return nil, err
	}
	if len(collectionNames) == 0 {
		connection.Logger.WithFields(fields).Warn("custom-mongo.base_connection.getCollection.ListCollectionNames collection not found")
	}

	collection := database.Collection(collectionName)

	baseCollection := collections.NewBaseCollection[T](collection, connection.Logger)

	return &baseCollection, nil
}
