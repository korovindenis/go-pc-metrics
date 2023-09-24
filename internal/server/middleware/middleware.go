package middleware

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"regexp"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type log interface {
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

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

func ErrorLogging(log log) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			log.Error("Error: ", zap.Error(c.Errors[0].Err))
		}
	}
}

func Gzip() gin.HandlerFunc {
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
					c.AbortWithError(http.StatusBadRequest, entity.ErrInvalidGzipData)
					return
				}
				requestBody = buf.Bytes()
			}
		}

		if !isGzip {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(safe); err != nil {
				c.AbortWithError(http.StatusBadRequest, entity.ErrReadingRequestBody)
				return
			}
			requestBody = buf.Bytes()
		}

		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = http.MaxBytesReader(c.Writer, io.NopCloser(bf), maxMemory)

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

func GzipResponse() gin.HandlerFunc {
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

func SetSign(secretKey, patternSign string) gin.HandlerFunc {
	return func(c *gin.Context) {

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, entity.ErrReadingRequestBody)
			return
		}
		if regexp.MustCompile(patternSign).MatchString(c.FullPath()) {
			serverHashSHA256, _ := computeHMAC(body, secretKey)
			c.Writer.Header().Set("HashSHA256", serverHashSHA256)
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}

func CheckSign(log log, secretKey, patternSign string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientHashSHA256 := c.GetHeader("HashSHA256")
		if clientHashSHA256 == "" {
			c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
			return
		}

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)

		serverHashSHA256, _ := computeHMAC(body, secretKey)

		if regexp.MustCompile(patternSign).MatchString(c.FullPath()) && clientHashSHA256 != serverHashSHA256 {
			log.Info("Client HashSHA256: " + clientHashSHA256)
			log.Info("Server HashSHA256: " + serverHashSHA256)
			log.Error("Check sign was failed")

			c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Next()
	}
}

func computeHMAC(input []byte, key string) (string, error) {
	keyBytes := []byte(key)

	h := hmac.New(sha256.New, keyBytes)

	_, err := h.Write(input)
	if err != nil {
		return "", err
	}

	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
