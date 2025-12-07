package gateway

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/berkedev13/bitaksi-gateway-service/internal/config"
	"github.com/gin-gonic/gin"
)

type Gateway struct {
	driverBaseURL    string
	passengerBaseURL string
	client           *http.Client
}

func NewGateway(cfg *config.Config) *Gateway {
	return &Gateway{
		driverBaseURL:    cfg.DriverBaseURL,
		passengerBaseURL: cfg.PassengerBaseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (g *Gateway) proxy(c *gin.Context, baseURL string) {
	req := c.Request

	targetURL := baseURL + req.URL.Path
	if req.URL.RawQuery != "" {
		targetURL += "?" + req.URL.RawQuery
	}

	var body io.Reader
	if req.Body != nil {
		buf, err := io.ReadAll(req.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			return
		}
		body = bytes.NewReader(buf)
	}

	outReq, err := http.NewRequestWithContext(req.Context(), req.Method, targetURL, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create proxy request"})
		return
	}

	for k, v := range req.Header {
		for _, vv := range v {
			outReq.Header.Add(k, vv)
		}
	}

	resp, err := g.client.Do(outReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unreachable"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read upstream response"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	c.Data(resp.StatusCode, contentType, respBody)
}

// ProxyDrivers godoc
// @Summary      Proxy driver-service requests
// @Description  Proxies /drivers isteklerini driver-service'e iletir
// @Tags         drivers, gateway
// @Produce      json
// @Param        Authorization  header    string  false  "Bearer token"
// @Success      200  {object}  map[string]any
// @Failure      401  {object}  map[string]any
// @Failure      429  {object}  map[string]any
// @Failure      502  {object}  map[string]any
// @Router       /drivers [get]
// @Router       /drivers [post]
func (g *Gateway) ProxyDrivers(c *gin.Context) {
	g.proxy(c, g.driverBaseURL)
}

// ProxyPassengers godoc
// @Summary      Proxy passenger-service requests
// @Description  Proxies /passengers isteklerini passenger-service'e iletir
// @Tags         passengers, gateway
// @Produce      json
// @Param        Authorization  header    string  false  "Bearer token"
// @Success      200  {object}  map[string]any
// @Failure      401  {object}  map[string]any
// @Failure      429  {object}  map[string]any
// @Failure      502  {object}  map[string]any
// @Router       /passengers [get]
// @Router       /passengers [post]
func (g *Gateway) ProxyPassengers(c *gin.Context) {
	g.proxy(c, g.passengerBaseURL)
}
