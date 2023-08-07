package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

func CheckMethodAndContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		namedURL := entity.ReqURI{
			MetricType: c.Param("metricType"),
			MetricName: c.Param("metricName"),
			MetricVal:  c.Param("metricVal"),
		}

		if c.Request.Method == http.MethodGet {
			if c.Request.URL.Path == "/" || (namedURL.MetricType != "" && namedURL.MetricName != "") {
				c.Next()
				return
			}

			c.AbortWithError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}

		if c.Request.Method == http.MethodPost {
			if namedURL.MetricType == "" || namedURL.MetricName == "" || namedURL.MetricVal == "" {
				c.AbortWithError(http.StatusBadRequest, errors.New("invalid URL format"))
				return
			}

			c.Next()
			return
		}

		c.AbortWithError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}
