package db

import (
	"context"
	"log"
	"time"

	"github.com/berkedev13/bitaksi-driver-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Client     *mongo.Client
	Database   *mongo.Database
	DriverColl *mongo.Collection
}

func NewMongoConnection(cfg *config.Config) *MongoConnection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("[mongo] connect error: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("[mongo] ping error: %v", err)
	}

	db := client.Database(cfg.MongoDatabase)
	driverColl := db.Collection(cfg.MongoDriverCollection)

	log.Println("[mongo] connected successfully")

	return &MongoConnection{
		Client:     client,
		Database:   db,
		DriverColl: driverColl,
	}
}
