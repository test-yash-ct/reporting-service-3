package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/healthops/reporting-service/internal/auth"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" || len(id) > 128 {
			b := make([]byte, 16)
			_, _ = rand.Read(b)
			id = hex.EncodeToString(b)
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("http method=%s path=%s status=%d dur_ms=%d", c.Request.Method, c.FullPath(), c.Writer.Status(), time.Since(start).Milliseconds())
	}
}

func Authenticate(secret string, maxTTLSec int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.ParseBearer(c.GetHeader("Authorization"), secret, time.Duration(maxTTLSec)*time.Second)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func ClaimsFrom(c *gin.Context) (auth.Claims, bool) {
	v, ok := c.Get("claims")
	if !ok {
		return auth.Claims{}, false
	}
	cl, ok := v.(auth.Claims)
	return cl, ok
}

func CORS(allow []string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, o := range allow {
		allowed[o] = struct{}{}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}
		if _, ok := allowed[origin]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "origin_denied"})
			return
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type windowLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int
	window time.Duration
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	l := &windowLimiter{hits: map[string][]time.Time{}, limit: limit, window: window}
	return func(c *gin.Context) {
		key := c.ClientIP()
		if cl, ok := ClaimsFrom(c); ok {
			key = cl.Tenant + ":" + cl.Sub
		}
		now := time.Now()
		l.mu.Lock()
		cut := now.Add(-window)
		arr := l.hits[key]
		kept := arr[:0]
		for _, t := range arr {
			if t.After(cut) {
				kept = append(kept, t)
			}
		}
		ok := len(kept) < limit
		if ok {
			l.hits[key] = append(kept, now)
		} else {
			l.hits[key] = kept
		}
		l.mu.Unlock()
		if !ok {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate_limited"})
			return
		}
		c.Next()
	}
}
