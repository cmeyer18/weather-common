package custom_mongo

import (
	"context"
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
func (bc *BaseConnection) GetCollection(database, collection string) (*mongo.Collection, error) {
	fields := logrus.Fields{"database": database, "collection": collection}

	// Alert if no db found
	dbNames, err := bc.client.ListDatabaseNames(context.TODO(), bson.M{"name": database})
	if err != nil {
		bc.Logger.WithError(err).Error("custom-mongo.base_connection.getCollection.ListDatabaseNames.error")
		return nil, err
	}

	if len(dbNames) == 0 {
		bc.Logger.WithFields(fields).Warn("custom-mongo.base_connection.getCollection.ListDatabasesNames db not found")
	}

	weatherDB := bc.client.Database(database)

	// Alert if no collection found
	collectionNames, err := weatherDB.ListCollectionNames(context.TODO(), bson.M{"name": collection})
	if err != nil {
		bc.Logger.WithError(err).Error("custom-mongo.base_connection.getCollection.ListCollectionNames.error")
		return nil, err
	}
	if len(collectionNames) == 0 {
		bc.Logger.WithFields(fields).Warn("custom-mongo.base_connection.getCollection.ListCollectionNames collection not found")
	}

	subscriberCollection := weatherDB.Collection(collection)
	return subscriberCollection, nil
}
