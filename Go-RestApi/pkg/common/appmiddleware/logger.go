package appmiddleware

import (
	"go-restapi/pkg/common/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		logger.Info("incoming request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrapped.statusCode),
			zap.Duration("duration", duration),
			zap.Int("bytes", wrapped.bytesWritten),
			zap.String("remote_addr", r.RemoteAddr),
		}

		if wrapped.statusCode >= 500 {
			logger.Error("request completed with server error", fields...)
		} else if wrapped.statusCode >= 400 {
			logger.Warn("request completed with client error", fields...)
		} else {
			logger.Info("request completed", fields...)
		}
	})
}
