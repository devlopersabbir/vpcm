package database

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	_ "modernc.org/sqlite"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database
var SQLiteDB *sql.DB

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

// InitSQLite establishes a connection to SQLite and pings the host.
func InitSQLite(path string) (*sql.DB, error) {
	slog.Debug("Connecting to SQLite", "path", path)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Without a busy timeout, a write that races another process (the background
	// API server, the desktop app) fails immediately with "database is locked"
	// instead of waiting for the other writer to finish.
	dsn := "file:" + path + "?_pragma=busy_timeout(5000)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	SQLiteDB = db
	return db, nil
}
