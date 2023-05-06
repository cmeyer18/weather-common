package custom_mongo

import (
	"context"
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

func NewConnection(credentials options.Credential, url string, client *mongo.Client) BaseConnection {
	return BaseConnection{credentials: credentials, mongoConnectionURL: url, client: client}
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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+bc.mongoConnectionURL).SetAuth(bc.credentials))
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
