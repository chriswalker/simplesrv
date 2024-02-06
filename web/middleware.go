package web

import "net/http"

// LogRequest logs details about incoming requests.
func (a *app) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("received request",
			"ip", r.RemoteAddr,
			"proto", r.Proto,
			"method", r.Method,
			"uri", r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}
