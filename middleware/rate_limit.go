package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Client holds the rate limiter instance and the last time it was active
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips sync.Map // Thread safe map to hold IP addresses
	r   rate.Limit
	b   int
}

// Configure a limiter with a rate and capacity
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		r: r,
		b: b,
	}

	// Background cleanup routine to get rid of dead tracking sessions
	go i.cleanupClients()

	return i
}

// Intercept the request and evaluate whether the IP address has exceeded its threshold
func (i *IPRateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract IP address of the client
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, `{"error": "Internal server error parsing client address}`, http.StatusInternalServerError)
			return
		}

		// Fetch or create rate limiter bucket for IP
		limiter := i.getLimiter(ip)

		// Evaluate if taking a token violates the limit
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Too many requests. Please slow down.}`))
			return
		}

		// Pass control if token was available
		next.ServeHTTP(w, r)
	}
}

func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	val, exists := i.ips.Load(ip)
	if exists {
		c := val.(*client)
		c.lastSeen = time.Now() // Update active timestamp
		return c.limiter
	}

	// Create new token bucket: r = refill rate, b = capacity
	limiter := rate.NewLimiter(i.r, i.b)

	i.ips.Store(ip, &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	})

	return limiter
}

func (i *IPRateLimiter) cleanupClients() {
	for {
		time.Sleep(1 * time.Minute)
		i.ips.Range(func(key, value interface{}) bool {
			c := value.(*client)
			// If client has not made request in 3 minutes, drop the client from memory tracking
			if time.Since(c.lastSeen) > 3*time.Minute {
				i.ips.Delete(key)
			}
			return true
		})
	}
}
