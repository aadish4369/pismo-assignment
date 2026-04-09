package middleware

import (
	"bytes"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxBodyLogLen = 2048

type responseBodyWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

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

		respBuf := bytes.NewBuffer(nil)
		c.Writer = &responseBodyWriter{ResponseWriter: c.Writer, buf: respBuf}

		c.Next()

		apiPath := c.FullPath()
		if apiPath == "" {
			apiPath = rawPath
		}
		respLine := "—"
		if b := respBuf.Bytes(); len(b) > 0 {
			s := string(b)
			if len(s) > maxBodyLogLen {
				s = s[:maxBodyLogLen] + "...(truncated)"
			}
			respLine = s
		}
		log.Printf("api: %s %s", c.Request.Method, apiPath)
		log.Printf("req: %s", reqLine)
		log.Printf("response: %d %s body: %s", c.Writer.Status(), time.Since(start).Round(time.Millisecond), respLine)
	}
}
