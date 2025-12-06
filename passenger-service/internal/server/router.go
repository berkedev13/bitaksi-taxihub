package server

import (
	"net/http"

	"github.com/berkedev13/bitaksi-passenger-service/internal/config"
	"github.com/berkedev13/bitaksi-passenger-service/internal/db"
	"github.com/berkedev13/bitaksi-passenger-service/internal/passenger"
	"github.com/gin-gonic/gin"

	_ "github.com/berkedev13/bitaksi-passenger-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	client, err := db.Connect(cfg.MongoURI)
	if err != nil {
		panic(err)
	}

	col := client.Database(cfg.DBName).Collection(cfg.PassengerCollection)

	repo := passenger.NewRepository(col)
	svc := passenger.NewService(repo)
	handler := passenger.NewHandler(svc)
	handler.RegisterRoutes(r)

	return r
}
