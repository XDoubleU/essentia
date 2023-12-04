package middleware

import (
	"github.com/XDoubleU/essentia/internal/core"
	"github.com/rs/cors"
)

func Cors() core.HandlerFunc {
	cors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
	})

	return func(c *core.Context) {
		cors.HandlerFunc(c.Writer.ResponseWriter, c.Request)
		c.Next()
	}
}
