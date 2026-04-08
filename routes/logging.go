package routes

import (
	"bytes"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxBodyLogLen = 2048

// APILogging logs one line each for api, req, and response. Skips /swagger entirely.
func APILogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}

		start := time.Now()
		rawPath := c.Request.URL.Path

		var reqLine string
		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewReader(body))
				switch {
				case len(body) > 0:
					s := string(body)
					if len(s) > maxBodyLogLen {
						s = s[:maxBodyLogLen] + "...(truncated)"
					}
					reqLine = s
				case c.Request.URL.RawQuery != "":
					reqLine = "?" + c.Request.URL.RawQuery
				}
			}
		}
		if reqLine == "" && c.Request.URL.RawQuery != "" {
			reqLine = "?" + c.Request.URL.RawQuery
		}
		if reqLine == "" {
			reqLine = "—"
		}

		c.Next()

		apiPath := c.FullPath()
		if apiPath == "" {
			apiPath = rawPath
		}
		log.Printf("api: %s %s", c.Request.Method, apiPath)
		log.Printf("req: %s", reqLine)
		log.Printf("response: %d %s", c.Writer.Status(), time.Since(start).Round(time.Millisecond))
	}
}
