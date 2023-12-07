package middleware

import (
	"github.com/XDoubleU/essentia/pkg/router"
	"github.com/rs/cors"
)

func Cors() router.HandlerFunc {
	cors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
	})

	return func(c *router.Context) {
		cors.HandlerFunc(c.Writer.ResponseWriter, c.Request)
		c.Next()
	}
}
