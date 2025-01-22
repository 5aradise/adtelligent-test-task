package middleware

import (
	"log/slog"
	"net/http"
)

func Recoverer(l *slog.Logger) Middleware {
	l.Info("recoverer middleware enabled")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						panic(rvr)
					}

					l.Error("recoverer middleware",
						slog.Any("error", rvr),
					)

					if r.Header.Get("Connection") != "Upgrade" {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
