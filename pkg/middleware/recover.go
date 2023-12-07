package middleware

import (
	"log"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/router"
)

func Recover() router.HandlerFunc {
	return func(c *router.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Writer.Header().Set("Connection", "close")
				c.Writer.WriteHeader(http.StatusInternalServerError)
				log.Printf("PANIC: %s\n", err)
			}
		}()

		c.Next()
	}
}
