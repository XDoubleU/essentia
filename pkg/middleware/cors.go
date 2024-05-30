package middleware

import (
	"github.com/rs/cors"
)

func Cors(allowedOrigins []string, useSentry bool) middleware {
	allowedHeaders := []string{"Content-Type"}
	if useSentry {
		allowedHeaders = append(allowedHeaders, "Baggage", "Sentry-Trace")
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   allowedHeaders,
	})

	return cors.Handler
}
