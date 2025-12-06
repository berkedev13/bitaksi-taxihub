// Package main TaxiHub Driver Service
//
// @title       TaxiHub Driver Service API
// @version     1.0
// @description Driver management microservice for Bitaksi TaxiHub case
// @host        localhost:8081
// @BasePath    /

package main

import (
	"log"

	"github.com/berkedev13/bitaksi-driver-service/internal/config"
	"github.com/berkedev13/bitaksi-driver-service/internal/server"
)

func main() {
	cfg := config.Load()
	r := server.NewRouter(cfg)

	log.Printf("[server] starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("[server] failed to start: %v", err)
	}
}
