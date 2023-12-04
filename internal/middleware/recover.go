package middleware

import (
	"log"
	"net/http"

	"github.com/XDoubleU/essentia/internal/core"
)

func Recover() core.HandlerFunc {
	return func(c *core.Context) {
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
