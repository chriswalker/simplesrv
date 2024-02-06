// Package web contains the main application server, which handles the
// server routes, handlers and associated HTML templates.
package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chriswalker/simplesrv/model"
	"github.com/chriswalker/simplesrv/service"
)

// HTTP routing and method checking taken from:
//
// https://github.com/benhoyt/go-routing/blob/9a2fa7a643ecb5681f504b95064d948ee2177c9a/retable/route.go
//
// It can be refactored out when the new routing changes come in with
// Go 1.22.

// itemService is an interface to the application's business logic service;
// used primarily for the web package's unit tests.
type itemService interface {
	GetItems() ([]model.Item, error)
}

// app is the basic application object, which stores business logic services,
// parsed templates and the app logger.
type app struct {
	svc       itemService
	routes    []route
	templates *template.Template
	logger    *slog.Logger
}

//go:embed static
var static embed.FS

// Start configures a new server object for handling requests and
// starts the HTTP server up.
func Start(addr, db string) error {
	app := &app{
		routes: make([]route, 0),
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	svc, err := service.NewItemService(db, app.logger)
	if err != nil {
		return err
	}
	app.svc = svc

	// Templates.
	t, err := app.loadTemplates()
	if err != nil {
		return err
	}
	app.templates = t

	mux := http.NewServeMux()

	// Set up static assets.
	fs, err := fs.Sub(static, "static")
	if err != nil {
		return err
	}
	fsrv := http.FileServer(http.FS(fs))
	mux.Handle("/static/", http.StripPrefix("/static/", fsrv))

	// Routes.
	app.routes = append(app.routes,
		newRoute(http.MethodGet, "/", app.index),
	)

	// Everything else handled by the app's Serve method.
	mux.HandleFunc("/", app.Serve)

	srv := http.Server{
		Addr:    addr,
		Handler: app.LogRequest(mux),
	}
	app.logger.Info(fmt.Sprintf("Listening on %s...", addr))
	return srv.ListenAndServe()
}

//go:embed templates
var templates embed.FS

// loadTemplates walks the templates directory, and creates templates from
// each file found there.
//
// Templates will be named after their filename.
func (a *app) loadTemplates() (*template.Template, error) {
	fs, err := fs.Sub(templates, "templates")
	if err != nil {
		a.logger.Error("unable to set up file system subtree: %v\n", err)
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"humanDate": func(t time.Time) string {
			f := "02 Jan 2006"
			sub := time.Since(t)
			if sub.Hours() < float64(24) {
				f = "15:04"
			}
			return t.Format(f)
		},
	}

	t, err := template.New("simplesrv").Funcs(funcMap).ParseFS(fs, "*.tmpl")
	if err != nil {
		return nil, err
	}

	return t, nil
}

//
// TODO
// Pretty much everything below can disappear once Go 1.22 lands, with
// various HTTP routing changes in it.
//

// route encapsulates a single route, specifying the HTTP method allowed,
// a regex to match against and a handler to invoke when matched.
type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

// newRoute creates and returns a configured route struct.
func newRoute(method, pattern string, handler http.HandlerFunc) route {
	route := route{
		method:  method,
		regex:   regexp.MustCompile("^" + pattern + "$"),
		handler: handler,
	}
	return route
}

type ctxKey struct{}

// Serve attempts to match the request path with registered handlers'
// regexes; on a match it then checks whether the method for that route is
// permitted.
//
// - No handler matches return 404 (Not Found) responses.
// - Invalid methods return 405 (Method Not Allowed) responses.
func (a *app) Serve(w http.ResponseWriter, r *http.Request) {
	var validMethods []string
	for _, route := range a.routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				validMethods = append(validMethods, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(validMethods) > 0 {
		w.Header().Set("Allow", strings.Join(validMethods, ","))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

// getField gets a path parameter from the request context, by
// slice index.
//nolint:unused
func getField(r *http.Request, idx int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[idx]
}

// formatResponseCode is a utility method to generate a nice string
// for a HTTP response code. Given a code such as 405, it will return
// "(405) Method Not Allowed".
//nolint:unused
func formatResponseCode(code int) string {
	return fmt.Sprintf("(%d) %s", code, http.StatusText(code))
}
