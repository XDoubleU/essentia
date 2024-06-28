package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimit(rps rate.Limit, bucketSize int, cleanupTimer time.Duration, removeAfter time.Duration) middleware {
	var (
		mu      sync.RWMutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(cleanupTimer)

			mu.RLock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > removeAfter {
					mu.RUnlock()
					mu.Lock()

					delete(clients, ip)

					mu.Unlock()
					mu.RLock()
				}
			}

			mu.RUnlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return rateLimit(&mu, clients, rps, bucketSize, next)
	}
}

func rateLimit(mu *sync.RWMutex, clients map[string]*client, rps rate.Limit, bucketSize int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			httptools.ServerErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rps, bucketSize)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			httptools.RateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
