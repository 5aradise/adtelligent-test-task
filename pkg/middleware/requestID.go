package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
)

var reqid atomic.Uint64

var RequestIDHeader = "X-Request-Id"

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = iota

func RequestID(l *slog.Logger) Middleware {
	l.Info("request id middleware enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				myid := reqid.Add(1)
				requestID = fmt.Sprintf("%015d", myid)
			}
			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestID(r *http.Request) string {
	ida := r.Context().Value(RequestIDKey)
	if ida == nil {
		return ""
	}

	id, ok := ida.(string)
	if !ok {
		return ""
	}
	return id
}
