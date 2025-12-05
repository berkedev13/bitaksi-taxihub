package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI              string
	MongoDatabase         string
	MongoDriverCollection string
	ServerPort            string
}

func Load() *Config {
	cfg := &Config{
		MongoURI:              getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:         getEnv("MONGO_DB", "bitaksi"),
		MongoDriverCollection: getEnv("MONGO_DRIVER_COLLECTION", "drivers"),
		ServerPort:            getEnv("PORT", "8081"),
	}

	log.Printf("[config] Loaded config: MongoURI=%s DB=%s Collection=%s Port=%s",
		cfg.MongoURI, cfg.MongoDatabase, cfg.MongoDriverCollection, cfg.ServerPort)

	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
