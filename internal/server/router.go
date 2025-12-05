package server

import (
	"github.com/berkedev13/bitaksi-driver-service/internal/config"
	"github.com/berkedev13/bitaksi-driver-service/internal/db"
	"github.com/berkedev13/bitaksi-driver-service/internal/driver"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	mongoConn := db.NewMongoConnection(cfg)

	repo := driver.NewRepository(mongoConn)
	service := driver.NewService(repo)
	handler := driver.NewHandler(service)

	handler.RegisterRoutes(r)

	return r
}
