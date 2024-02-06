package web

import (
	"net/http"
)

// serverError logs an error and returns an Internal Server Error response.
func (a *app) serverError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error(err.Error(), "uri", r.URL.RequestURI(), "method", r.Method)
	clientError(w, http.StatusInternalServerError)
}

// clientError returns the supplied error status to the client.
func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound returns a Not Found error response.
//nolint:unused
func (a *app) notFoundError(w http.ResponseWriter) {
	clientError(w, http.StatusNotFound)
}
