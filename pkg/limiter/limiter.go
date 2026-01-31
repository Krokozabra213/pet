package limiter

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	lastSeen time.Time
	limiter  *rate.Limiter
}

func NewVisitor(t time.Time, limiter *rate.Limiter) *visitor {
	return &visitor{
		lastSeen: t,
		limiter:  limiter,
	}
}

type rateLimiter struct {
	mu sync.RWMutex

	visitors map[string]*visitor
	limit    rate.Limit
	burst    int
	ttl      time.Duration
}

func newRateLimiter(limit, burst int, ttl time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:    rate.Limit(limit),
		burst:    burst,
		ttl:      ttl,
		visitors: make(map[string]*visitor),
	}
}

func (r *rateLimiter) getLimiter(ip string) *rate.Limiter {
	r.mu.RLock()
	visitor, exists := r.visitors[ip]
	r.mu.RUnlock()

	if !exists {
		limiter := rate.NewLimiter(r.limit, r.burst)
		visitor = NewVisitor(time.Now(), limiter)
		r.mu.Lock()
		r.visitors[ip] = visitor
		r.mu.Unlock()
		return limiter
	}

	visitor.lastSeen = time.Now()
	return visitor.limiter
}

func (r *rateLimiter) cleanVisitors() {
	for {
		r.mu.Lock()
		for ip, v := range r.visitors {
			if time.Since(v.lastSeen) > r.ttl {
				delete(r.visitors, ip)
			}
		}
		r.mu.Unlock()
		time.Sleep(90 * time.Second)
	}
}

func Limit(rps, burst int, ttl time.Duration) gin.HandlerFunc {
	l := newRateLimiter(rps, burst, ttl)

	go l.cleanVisitors()

	return func(c *gin.Context) {
		var err error
		ip := c.Request.RemoteAddr
		ip, _, err = net.SplitHostPort(ip)
		if err != nil {
			slog.Error("failed parse client ip address", "err", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if !l.getLimiter(ip).Allow() {
			slog.Info("attempt ddos atack...")
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}
