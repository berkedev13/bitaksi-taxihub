package config

import (
	"os"
)

type Config struct {
	Port                string
	MongoURI            string
	DBName              string
	PassengerCollection string
}

func Load() (*Config, error) {
	port := os.Getenv("PASSENGER_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "bitaksi"
	}

	passengerCol := os.Getenv("PASSENGER_COLLECTION")
	if passengerCol == "" {
		passengerCol = "passengers"
	}

	return &Config{
		Port:                port,
		MongoURI:            mongoURI,
		DBName:              dbName,
		PassengerCollection: passengerCol,
	}, nil
}
