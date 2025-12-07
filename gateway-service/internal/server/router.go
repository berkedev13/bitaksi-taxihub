package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/berkedev13/bitaksi-gateway-service/internal/config"
	"github.com/berkedev13/bitaksi-gateway-service/internal/gateway"
	"github.com/berkedev13/bitaksi-gateway-service/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Health godoc
// @Summary      Gateway health check
// @Description  Returns gateway service status
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "gateway",
	})
}

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Global rate limit: 5 req/s, burst 10
	clientLimiter := middleware.NewClientLimiter(rate.Limit(5), 10)
	r.Use(middleware.RateLimitMiddleware(clientLimiter))

	// Public endpoints (JWT yok)
	r.GET("/health", healthHandler)

	// Swagger UI (public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected endpoints (JWT zorunlu)
	api := r.Group("/")
	api.Use(middleware.JWTAuthMiddleware())
	api.Use(middleware.APIKeyMiddleware())

	gw := gateway.NewGateway(cfg)

	// /drivers -> driver-service
	api.Any("/drivers", gw.ProxyDrivers)
	api.Any("/drivers/*any", gw.ProxyDrivers)

	// /passengers -> passenger-service
	api.Any("/passengers", gw.ProxyPassengers)
	api.Any("/passengers/*any", gw.ProxyPassengers)

	return r
}
