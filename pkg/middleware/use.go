package middleware

import "net/http"

func Use(router http.Handler, mids ...Middleware) http.Handler {
	for i := len(mids) - 1; i >= 0; i-- {
		router = mids[i](router)
	}
	return router
}
