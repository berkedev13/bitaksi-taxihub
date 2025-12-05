package driver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	drivers := r.Group("/drivers")
	{
		drivers.POST("", h.CreateDriver)
		drivers.PUT("/:id", h.UpdateDriver)
	}
}

func (h *Handler) CreateDriver(c *gin.Context) {
	var req CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	driver, err := h.service.CreateDriver(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create driver",
		})
		return
	}

	c.JSON(http.StatusCreated, driver)
}

func (h *Handler) UpdateDriver(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "UpdateDriver is not implemented yet (Day 2)",
	})
}
