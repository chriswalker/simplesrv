package web

import (
	"net/http"
)

// index generates the index page, retrieving stored items from
// the databse for display on the page.
func (a *app) index(w http.ResponseWriter, r *http.Request) {
	items, err := a.svc.GetItems()
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.renderTemplate(w, r, http.StatusOK, "index.tmpl", items)
}
