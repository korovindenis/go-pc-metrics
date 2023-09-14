package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
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
		c.Next()
	}
}

func GzipMiddleware() gin.HandlerFunc {
	var maxMemory int64 = 64 << 20 // 64 MB

	return func(c *gin.Context) {
		var requestBody []byte
		isGzip := false
		safe := &io.LimitedReader{R: c.Request.Body, N: maxMemory}

		if c.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(safe)
			if err == nil {
				isGzip = true
				var buf bytes.Buffer
				if _, err := buf.ReadFrom(reader); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Gzip data"})
					c.Abort()
					return
				}
				requestBody = buf.Bytes()
			}
		}

		if !isGzip {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(safe); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
				c.Abort()
				return
			}
			requestBody = buf.Bytes()
		}

		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = http.MaxBytesReader(c.Writer, ioutil.NopCloser(bf), maxMemory)

		c.Next()
	}
}

type gzipResponseWriter struct {
	io.Writer
	gin.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Writer.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(c.Writer)
			defer gzipWriter.Close()

			c.Writer = &gzipResponseWriter{Writer: gzipWriter, ResponseWriter: c.Writer}
		}
		c.Next()
	}
}
