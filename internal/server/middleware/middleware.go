package middleware

import (
	"compress/gzip"
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

			// body, err := ioutil.ReadAll(c.Request.Body)
			// if err != nil {
			// 	c.JSON(500, gin.H{"error": "Ошибка чтения тела запроса"})
			// 	return
			// }
			// fmt.Printf("Содержимое запроса: %s\n", string(body))

		} else {
			// body, err := ioutil.ReadAll(c.Request.Body)
			// if err != nil {
			// 	c.JSON(500, gin.H{"error": "Ошибка чтения тела запроса"})
			// 	return
			// }
			// fmt.Printf("Содержимое запроса: %s\n", string(body))

		}

		c.Next()
	}
}
