# 🚕 TaxiHub -- Go Microservices Case Study

**Driver-Service • Passenger-Service • API Gateway • MongoDB • JWT •
Rate Limit • Swagger • Docker Compose**

This project is a modern microservice-based backend system designed for
the TaxiHub case. It includes separate Driver and Passenger services, a
centralized API Gateway, JWT authentication, API Key validation, rate
limiting, Swagger documentation, and full Docker Compose orchestration.

------------------------------------------------------------------------

## 📦 Architecture

    bitaksi-taxihub/
    │
    ├── driver-service/        # Driver CRUD + Nearby Search + MongoDB
    ├── passenger-service/     # Passenger CRUD + Nearby Search + MongoDB
    ├── gateway-service/       # API Gateway (Proxy + JWT + Rate Limit + API Key)
    │
    └── docker-compose.yml     # Starts all services together

### 🔧 Services Overview

  --------------------------------------------------------------------------
  Service                 Description                      Port
  ----------------------- -------------------------------- -----------------
  **API Gateway**         Central entry point, proxy       `8080`
                          routing, JWT + API Key + rate    
                          limit                            

  **Driver-Service**      Driver CRUD, pagination,         `8081`
                          geolocation, MongoDB             

  **Passenger-Service**   Passenger CRUD, pagination,      `8082`
                          geolocation, MongoDB             

  **MongoDB**             Shared database                  `27017`
  --------------------------------------------------------------------------

------------------------------------------------------------------------

## 🚀 Features

### ✔ Driver-Service

-   Create / Update drivers
-   Paginated driver listing
-   Find nearby drivers (6 km)
-   MongoDB repository pattern
-   Swagger/OpenAPI docs
-   Unit tests (service layer)
-   Clean architecture (handler → service → repository)

### ✔ Passenger-Service

-   Create / Update passengers
-   Paginated passenger listing
-   Find nearby passengers (6 km)
-   MongoDB repository pattern
-   Swagger/OpenAPI docs
-   Unit tests

### ✔ API Gateway

-   Routes all requests to driver/passenger services
-   JWT Authentication
-   API Key validation
-   Rate Limiting (5 requests/sec, burst 10)
-   Swagger (optional)
-   Middleware unit tests

### ✔ Bonus Features (All Completed)

-   Docker Compose environment
-   JWT Auth
-   API Key middleware
-   Rate Limiting middleware
-   Swagger/OpenAPI
-   Unit tests

------------------------------------------------------------------------

## 🐳 Running with Docker

Start the entire system:

``` bash
docker compose up --build
```

All services will start automatically.

------------------------------------------------------------------------

## 🔗 Swagger Documentation

  Service                 URL
  ----------------------- ------------------------------------------
  **Driver-Service**      http://localhost:8081/swagger/index.html
  **Passenger-Service**   http://localhost:8082/swagger/index.html
  **Gateway-Service**     http://localhost:8080/swagger/index.html

------------------------------------------------------------------------

## 🔐 Authentication

Gateway requires **both JWT Token and API Key** for all `/drivers` and
`/passengers` requests.

### ✔ Set JWT Token

``` bash
export TOKEN="your_jwt_token_here"
```

### ✔ Set API Key

``` bash
export API_KEY="apikey123"
```

------------------------------------------------------------------------

## 🧪 Example Requests

### Create Driver

``` bash
curl -X POST http://localhost:8080/drivers   -H "Authorization: Bearer $TOKEN"   -H "X-API-Key: $API_KEY"   -H "Content-Type: application/json"   -d '{
    "firstName": "Ahmet",
    "lastLastName": "Demir",
    "plate": "34ABC123",
    "taxiType": "sari",
    "carBrand": "Toyota",
    "carModel": "Corolla",
    "lat": 41.0431,
    "lon": 29.0099
  }'
```

### List Drivers

``` bash
curl "http://localhost:8080/drivers?page=1&pageSize=20"   -H "Authorization: Bearer $TOKEN"   -H "X-API-Key: $API_KEY"
```

### Find Nearby Drivers

``` bash
curl "http://localhost:8080/drivers/nearby?lat=41.0431&lon=29.0099"   -H "Authorization: Bearer $TOKEN"   -H "X-API-Key: $API_KEY"
```

------------------------------------------------------------------------

## 🧪 Running Unit Tests

### All tests in a service

``` bash
go test ./...
```

### Only gateway middleware tests

``` bash
go test ./internal/middleware
```

------------------------------------------------------------------------

## 📘 Technologies Used

-   **Go 1.23**
-   **Gin Web Framework**
-   **MongoDB**
-   **JWT Auth (HS256)**
-   **API Key Middleware**
-   **Rate Limiting (golang.org/x/time/rate)**
-   **Swagger (swaggo/swag)**
-   **Docker & Docker Compose**
-   **Unit Testing (Go testing)**

------------------------------------------------------------------------

## 🎯 Summary

This project demonstrates a fully functional, production-style
microservice architecture written in Go.\
It includes authentication, routing, proxying, caching, documentation,
containerization, and solid code structure.

Ideal for real-world case studies, portfolio projects, and interview
challenges.
