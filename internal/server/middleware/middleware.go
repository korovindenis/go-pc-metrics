package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// check GET or POST
func CheckMethod() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodGet {
			c.AbortWithError(http.StatusMethodNotAllowed, entity.ErrMethodNotAllowed)
			return
		}
		contentEncoding := c.GetHeader("Content-Encoding")
		if strings.ToLower(contentEncoding) == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				fmt.Println(err)
				c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
				return
			}
			defer reader.Close()

			var buf strings.Builder
			if _, err := io.Copy(&buf, reader); err != nil {
				c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
				return
			}

			c.Request.Body = io.NopCloser(strings.NewReader(buf.String()))
		}
		c.Next()
	}
}
