package config

import "os"

type Config struct {
	Port             string
	DriverBaseURL    string
	PassengerBaseURL string
}

func Load() *Config {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	driverURL := os.Getenv("DRIVER_SERVICE_URL")
	if driverURL == "" {
		driverURL = "http://localhost:8081"
	}

	passengerURL := os.Getenv("PASSENGER_SERVICE_URL")
	if passengerURL == "" {
		passengerURL = "http://localhost:8082"
	}

	return &Config{
		Port:             port,
		DriverBaseURL:    driverURL,
		PassengerBaseURL: passengerURL,
	}
}
