package middleware

import (
	"log"
	"time"

	"github.com/XDoubleU/essentia/pkg/router"
)

func Logger() router.HandlerFunc {
	return func(c *router.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.Request.Response.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
