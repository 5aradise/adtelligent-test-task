package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type WriterWrapper struct {
	http.ResponseWriter
	status int
}

func (ww *WriterWrapper) WriteHeader(statusCode int) {
	ww.status = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

func Logger(l *slog.Logger) Middleware {
	l.Info("logger middleware enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := &WriterWrapper{w, 0}
			beginReq := time.Now()
			next.ServeHTTP(ww, r)
			duration := time.Since(beginReq)

			l.Info("request info",
				slog.Any("status", ww.status),
				slog.Duration("duration", duration),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("id", GetRequestID(r)),
			)
		})
	}
}
