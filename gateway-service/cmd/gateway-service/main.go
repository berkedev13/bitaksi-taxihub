// Package main TaxiHub Gateway Service
//
// @title       TaxiHub Gateway API
// @version     1.0
// @description API Gateway for Bitaksi TaxiHub case
// @host        localhost:8080
// @BasePath    /
package main

import (
	"log"

	"github.com/berkedev13/bitaksi-gateway-service/internal/config"
	"github.com/berkedev13/bitaksi-gateway-service/internal/server"

	_ "github.com/berkedev13/bitaksi-gateway-service/docs"
)

func main() {
	cfg := config.Load()

	r := server.NewRouter(cfg)

	log.Printf("Gateway service is running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}
