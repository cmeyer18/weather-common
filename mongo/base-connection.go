package mongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection details a connection struct for the main file.
type Connection struct {
	Credentials        options.Credential
	MongoConnectionURL string
	Ctx                context.Context
	Client             *mongo.Client
}

// Ping tries to connect the server.
func (c *Connection) ping() bool {
	if err := c.Client.Ping(c.Ctx, readpref.Primary()); err != nil {
		log.WithError(err).Error("mongo.base_connection.ping.error")
		return false
	}
	log.Info("mongo.base_connection.ping.success")
	return true
}

// Disconnect takes the given ctx and client and disconnects them.
func (c *Connection) Disconnect() {
	err := c.Client.Disconnect(c.Ctx)
	if err != nil {
		log.WithError(err).Error("mongo.base_connection.disconnectClient.error")
	}
	log.Info("mongo.base_connection.disconnectClient.success")
}

// Connect sets up a connection to a mongoDB server.
func (c *Connection) Connect() error {
	client, err := mongo.Connect(c.Ctx, options.Client().ApplyURI("mongodb://"+c.MongoConnectionURL).SetAuth(c.Credentials))
	if err != nil {
		log.WithError(err).Error("mongo.base_connection.connectedClient.Connect.error")
		return err
	}
	c.Client = client

	// Ping for up status
	success := c.ping()
	if !success {
		log.Error("mongo.base_connection.connectedClient.Ping.error")
		c.Disconnect()
		return nil
	}

	log.Info("mongo.base_connection.connectedClient.Connect.success")
	return nil
}
