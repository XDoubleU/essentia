package middleware

import (
	"log"
	"time"

	"github.com/XDoubleU/essentia/internal/core"
)

func Logger() core.HandlerFunc {
	return func(c *core.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.Request.Response.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
