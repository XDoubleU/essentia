package middleware

import (
	"net"
	"sync"
	"time"

	"github.com/XDoubleU/essentia/pkg/router"
	"golang.org/x/time/rate"
)

func RateLimit() router.HandlerFunc {
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

	return func(c *router.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			//TODO: app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rps, bucketSize)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			//TODO: app.rateLimitExceededResponse(w, r)
			return
		}

		c.Next()
	}
}
