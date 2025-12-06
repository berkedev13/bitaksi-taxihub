package passenger

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	p := r.Group("/passengers")
	{
		p.POST("", h.Create)
		p.PUT("/:id", h.Update)
		p.GET("", h.List)
		p.GET("/nearby", h.GetNearby)
	}
}

// Create passenger godoc
// @Summary      Create a new passenger
// @Description  Creates a new passenger with basic info and current location
// @Tags         passengers
// @Accept       json
// @Produce      json
// @Param        passenger  body      CreatePassengerRequest  true  "Passenger info"
// @Success      201        {object}  Passenger
// @Failure      400        {object}  map[string]any
// @Failure      500        {object}  map[string]any
// @Router       /passengers [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreatePassengerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.svc.Create(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Update passenger godoc
// @Summary      Update a passenger
// @Description  Partially updates passenger info
// @Tags         passengers
// @Accept       json
// @Produce      json
// @Param        id         path      string                 true  "Passenger ID"
// @Param        passenger  body      UpdatePassengerRequest true  "Passenger update info"
// @Success      200        {object}  Passenger
// @Failure      400        {object}  map[string]any
// @Failure      404        {object}  map[string]any
// @Failure      500        {object}  map[string]any
// @Router       /passengers/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdatePassengerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.svc.Update(ctx, id, req)
	if err == ErrPassengerNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "passenger not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// List passengers godoc
// @Summary      List passengers
// @Description  Lists passengers with pagination
// @Tags         passengers
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        pageSize  query     int  false  "Page size"    default(20)
// @Success      200       {array}   Passenger
// @Failure      400       {object}  map[string]any
// @Failure      500       {object}  map[string]any
// @Router       /passengers [get]
func (h *Handler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("pageSize", "20")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.svc.List(ctx, int64(page), int64(size))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get nearby passengers godoc
// @Summary      Get nearby passengers
// @Description  Returns passengers within 6km radius of given location
// @Tags         passengers
// @Accept       json
// @Produce      json
// @Param        lat  query     number  true  "Latitude"
// @Param        lon  query     number  true  "Longitude"
// @Success      200  {array}   Passenger
// @Failure      400  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /passengers/nearby [get]
func (h *Handler) GetNearby(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

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

	res, err := h.svc.GetNearby(ctx, lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find nearby"})
		return
	}

	c.JSON(http.StatusOK, res)
}
