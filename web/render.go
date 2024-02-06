package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// renderTemplate gets the supplied template and executes it with the
// given data, emitting the results to the provided http.ResponseWriter.
func (a *app) renderTemplate(w http.ResponseWriter, r *http.Request,
	status int, tmpl string, data any) {
	t := a.templates.Lookup(tmpl)
	if t == nil {
		a.serverError(w, r, fmt.Errorf("template '%s' not found", tmpl))
		return
	}

	b := new(bytes.Buffer)
	if err := t.Execute(b, data); err != nil {
		a.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	_, _ = b.WriteTo(w)
}

// renderJSON marshals the supplied data to JSON, emitting the results to
// the provided http.ResponseWriter.
//nolint:unused
func (a *app) renderJSON(w http.ResponseWriter, r *http.Request,
	status int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	fmt.Fprint(w, string(b))
}
