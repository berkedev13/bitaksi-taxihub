// Package main TaxiHub Passenger Service
//
// @title       TaxiHub Passenger Service API
// @version     1.0
// @description Passenger management microservice for Bitaksi TaxiHub case
// @host        localhost:8082
// @BasePath    /
package main

import (
	"log"

	"github.com/berkedev13/bitaksi-passenger-service/internal/config"
	"github.com/berkedev13/bitaksi-passenger-service/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	r := server.NewRouter(cfg)

	log.Printf("Passenger service is running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
