// Package web contains the main application server, which handles the
// server routes, handlers and associated HTML templates.
package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/chriswalker/simplesrv/model"
	"github.com/chriswalker/simplesrv/service"
)

// itemService is an interface to the application's business logic service;
// used primarily for the web package's unit tests.
type itemService interface {
	GetItems() ([]model.Item, error)
}

// app is the basic application object, which stores business logic services,
// parsed templates and the app logger.
type app struct {
	svc       itemService
	templates *template.Template
	logger    *slog.Logger
}

//go:embed static
var static embed.FS

// Start configures a new server object for handling requests and
// starts the HTTP server up.
func Start(addr, db string) error {
	app := &app{
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
	mux.HandleFunc("GET /{$}", app.index)

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

// formatResponseCode is a utility method to generate a nice string
// for a HTTP response code. Given a code such as 405, it will return
// "(405) Method Not Allowed".
//nolint:unused
func formatResponseCode(code int) string {
	return fmt.Sprintf("(%d) %s", code, http.StatusText(code))
}
