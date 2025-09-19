package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holds the database connection
type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewConnection creates a new MongoDB connection
func NewConnection(uri, dbName string) (*DB, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Successfully connected to MongoDB database: %s", dbName)

	database := client.Database(dbName)

	return &DB{
		Client:   client,
		Database: database,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("MongoDB connection closed")
	return nil
}

// GetCollection returns a MongoDB collection
func (db *DB) GetCollection(name string) *mongo.Collection {
	return db.Database.Collection(name)
}
