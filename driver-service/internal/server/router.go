package server

import (
	"github.com/berkedev13/bitaksi-driver-service/internal/config"
	"github.com/berkedev13/bitaksi-driver-service/internal/db"
	"github.com/berkedev13/bitaksi-driver-service/internal/driver"
	"github.com/gin-gonic/gin"

	_ "github.com/berkedev13/bitaksi-driver-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	mongoConn := db.NewMongoConnection(cfg)

	repo := driver.NewRepository(mongoConn)
	service := driver.NewService(repo)
	handler := driver.NewHandler(service)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.RegisterRoutes(r)

	return r
}
