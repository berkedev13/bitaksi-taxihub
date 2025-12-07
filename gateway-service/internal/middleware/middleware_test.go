package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/berkedev13/bitaksi-gateway-service/internal/middleware"
	"golang.org/x/time/rate"
)

func createTestJWT(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "1234567890",
		"name":  "Test User",
		"admin": true,
		"iat":   time.Now().Unix(),
	})

	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "jwtsecret123")
	defer os.Unsetenv("JWT_SECRET")

	r := gin.New()
	r.Use(middleware.JWTAuthMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := createTestJWT("jwtsecret123")

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "jwtsecret123")
	defer os.Unsetenv("JWT_SECRET")

	r := gin.New()
	r.Use(middleware.JWTAuthMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer this_is_not_a_valid_token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "jwtsecret123")
	defer os.Unsetenv("JWT_SECRET")

	r := gin.New()
	r.Use(middleware.JWTAuthMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPIKeyMiddleware_ValidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("API_KEY", "apikey123")
	defer os.Unsetenv("API_KEY")

	r := gin.New()
	r.Use(middleware.APIKeyMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "apikey123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPIKeyMiddleware_InvalidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("API_KEY", "apikey123")
	defer os.Unsetenv("API_KEY")

	r := gin.New()
	r.Use(middleware.APIKeyMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-API-Key", "wrong-key")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAPIKeyMiddleware_MissingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("API_KEY", "apikey123")
	defer os.Unsetenv("API_KEY")

	r := gin.New()
	r.Use(middleware.APIKeyMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := middleware.NewClientLimiter(rate.Limit(1), 1)

	r := gin.New()
	r.Use(middleware.RateLimitMiddleware(limiter))

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}
