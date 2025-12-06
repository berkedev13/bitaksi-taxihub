package driver

import (
	"context"
	"net/http"
	"strconv"
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
		drivers.GET("", h.ListDrivers)
		drivers.GET("/nearby", h.GetNearbyDrivers)
	}
}

// CreateDriver godoc
// @Summary      Create a new driver
// @Description  Creates a new driver with basic information and current location
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        driver  body      CreateDriverRequest  true  "Driver info"
// @Success      201     {object}  Driver
// @Failure      400     {object}  map[string]any
// @Failure      500     {object}  map[string]any
// @Router       /drivers [post]
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
	id := c.Param("id")

	var req UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updated, err := h.service.UpdateDriver(ctx, id, req)
	if err != nil {
		if err == ErrDriverNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "driver not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update driver",
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// ListDrivers godoc
// @Summary      List drivers
// @Description  Lists drivers with pagination
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "Page number"   default(1)
// @Param        pageSize  query     int  false  "Page size"     default(20)
// @Success      200       {array}   Driver
// @Failure      400       {object}  map[string]any
// @Failure      500       {object}  map[string]any
// @Router       /drivers [get]
func (h *Handler) ListDrivers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageSize"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	drivers, err := h.service.ListDrivers(ctx, int64(page), int64(pageSize))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list drivers",
		})
		return
	}

	c.JSON(http.StatusOK, drivers)
}

// GetNearbyDrivers godoc
// @Summary      Get nearby drivers
// @Description  Returns drivers within 6km radius of given location and optional taxi type
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        lat       query     number  true   "Latitude"
// @Param        lon       query     number  true   "Longitude"
// @Param        taxiType  query     string  false  "Taxi type filter (sari, turkuaz, vb.)"
// @Success      200       {array}   NearbyDriver
// @Failure      400       {object}  map[string]any
// @Failure      500       {object}  map[string]any
// @Router       /drivers/nearby [get]
func (h *Handler) GetNearbyDrivers(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	taxiType := c.Query("taxiType")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lon"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	nearby, err := h.service.GetNearbyDrivers(ctx, lat, lon, taxiType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get nearby drivers",
		})
		return
	}

	c.JSON(http.StatusOK, nearby)
}
