package middleware

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
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
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithError(http.StatusMethodNotAllowed, entity.ErrMethodNotAllowed)
				return
			}

			reader, err := gzip.NewReader(bytes.NewReader(body))
			if err != nil {
				c.AbortWithError(http.StatusMethodNotAllowed, entity.ErrMethodNotAllowed)
				return
			}

			uncompressedBody, err := ioutil.ReadAll(reader)
			if err != nil {
				c.AbortWithError(http.StatusMethodNotAllowed, entity.ErrMethodNotAllowed)
				return
			}

			c.Request.Body = ioutil.NopCloser(bytes.NewReader(uncompressedBody))
		}

		c.Next()
	}
}
