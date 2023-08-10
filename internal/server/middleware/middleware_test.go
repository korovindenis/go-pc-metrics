package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/korovindenis/go-pc-metrics/internal/server/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCheckMethod(t *testing.T) {
	// Setup
	r := gin.New()
	r.Use(middleware.CheckMethod())

	// Define a handler for testing
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "GET request successful")
	})
	r.POST("/update/", func(c *gin.Context) {
		c.String(http.StatusOK, "POST request successful")
	})

	// Test GET request
	getReq := httptest.NewRequest(http.MethodGet, "/", nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)
	assert.Equal(t, http.StatusOK, getResp.Code)

	// Test POST request
	postReq := httptest.NewRequest(http.MethodPost, "/update", nil)
	postResp := httptest.NewRecorder()
	r.ServeHTTP(postResp, postReq)
	assert.Equal(t, http.StatusOK, postResp.Code)

	// Test unsupported method
	deleteReq := httptest.NewRequest(http.MethodDelete, "/", nil)
	deleteResp := httptest.NewRecorder()
	r.ServeHTTP(deleteResp, deleteReq)
	assert.Equal(t, http.StatusMethodNotAllowed, deleteResp.Code)
	assert.Contains(t, deleteResp.Body.String(), entity.ErrMethodNotAllowed)
}
