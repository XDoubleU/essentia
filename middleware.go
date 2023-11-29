package essentia

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/goddtriffin/helmet"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

func Helmet() HandlerFunc {
	helmet := helmet.Default()

	return func(c *Context) {
		helmet.ContentSecurityPolicy.Header(c.Writer.ResponseWriter)
		helmet.XContentTypeOptions.Header(c.Writer.ResponseWriter)
		helmet.XDNSPrefetchControl.Header(c.Writer.ResponseWriter)
		helmet.XDownloadOptions.Header(c.Writer.ResponseWriter)
		helmet.ExpectCT.Header(c.Writer.ResponseWriter)
		helmet.FeaturePolicy.Header(c.Writer.ResponseWriter)
		helmet.XFrameOptions.Header(c.Writer.ResponseWriter)
		helmet.XPermittedCrossDomainPolicies.Header(c.Writer.ResponseWriter)
		helmet.XPoweredBy.Header(c.Writer.ResponseWriter)
		helmet.ReferrerPolicy.Header(c.Writer.ResponseWriter)
		helmet.StrictTransportSecurity.Header(c.Writer.ResponseWriter)
		helmet.XXSSProtection.Header(c.Writer.ResponseWriter)

		c.Next()
	}
}

func Cors() HandlerFunc {
	cors := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
	})

	return func(c *Context) {
		cors.HandlerFunc(c.Writer.ResponseWriter, c.Request)
		c.Next()
	}
}

func Recover() HandlerFunc {
	return func(c *Context) {
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

func RateLimit() HandlerFunc {
	var rps rate.Limit = 10
	var bucketSize = 30

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(c *Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			//TODO: app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rps, bucketSize)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			//TODO: app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		c.Next()
	}
}
