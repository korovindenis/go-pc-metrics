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
	// return func(c *gin.Context) {
	// 	if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
	// 		// reader, err := gzip.NewReader(c.Request.Body)
	// 		// if err != nil {
	// 		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Gzip data"})
	// 		// 	c.Abort()
	// 		// 	return
	// 		// }
	// 		// defer reader.Close()
	// 		//

	// 		var requestBody []byte
	// 		var maxMemory int64 = 64 << 20 // 64 MB
	// 		reader, _ := gzip.NewReader(&io.LimitedReader{R: c.Request.Body, N: maxMemory})
	// 		requestBody, _ = ioutil.ReadAll(reader)
	// 		c.Request.Body.Close()
	// 		bf := bytes.NewBuffer(requestBody)
	// 		c.Request.Body = http.MaxBytesReader(c.Writer, ioutil.NopCloser(bf), maxMemory)
	// 		//
	// 		//c.Request.Body = http.MaxBytesReader(c.Writer, reader, c.Request.ContentLength)
	// 	}
	// 	c.Next()
	// }
	var maxMemory int64 = 64 << 20 // 64 MB

	return func(c *gin.Context) {
		var requestBody []byte
		isGzip := false
		safe := &io.LimitedReader{R: c.Request.Body, N: maxMemory}

		if c.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(safe)
			if err == nil {
				isGzip = true
				requestBody, _ = ioutil.ReadAll(reader)
			}
		}

		if !isGzip {
			requestBody, _ = ioutil.ReadAll(safe)
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
