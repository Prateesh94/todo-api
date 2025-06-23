package encrypt

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var clients = make(map[string]*Client)
var mu sync.Mutex

type Client struct {
	limiter *rate.Limiter
}

func getclient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	if client, exists := clients[ip]; exists {
		return client.limiter
	}
	limit := rate.NewLimiter(rate.Limit(50), 1)
	clients[ip] = &Client{limit}
	return limit
}

func Limitmid(a http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		lim := getclient(ip)
		if !lim.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		a.ServeHTTP(w, r)
	})
}
