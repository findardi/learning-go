package appmiddleware

import (
	"go-restapi/pkg/common"
	"go-restapi/pkg/config/limiter"
	"net"
	"net/http"
)

func RatelimiterMiddleware(limiter limiter.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.Header.Get("X-Real-IP")
			if ip == "" {
				ip = r.Header.Get("X-Forwarded-For")
			}
			if ip == "" {
				ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			}

			if allow, retryAfter := limiter.Allow(ip); !allow {
				common.RespondExcededError(w, r, retryAfter.String())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
