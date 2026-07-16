package database

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// InitMongo establishes a connection to MongoDB and ping the host.
func InitMongo(uri string, dbName string) (*mongo.Database, error) {
	slog.Debug("Connecting to MongoDB", "uri", uri, "db", dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	MongoClient = client
	MongoDB = client.Database(dbName)

	return MongoDB, nil
}
