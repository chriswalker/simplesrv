/*
simplesrv is a Go project template that provides a basic HTTP server
backing onto a simple database (in this case an in-process SQLite). It
provides enough structure for larger projects, including basic handlers,
a service and simple calls to the database.

Usage:

   simplesrv --sqlitedb ./app.db

Flags:

   -a, --addr
       Address/port for the server to listen on.
   --sqlitedb
       Path to the SQLite database file. It will be created if it does not exist.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chriswalker/simplesrv/web"
)

func main() {
	// Slightly nicer help output.
	flag.Usage = func() {
		fmt.Println("simplsrv - a Go project template for a minimalist server")
		fmt.Println("\nUsage: simplesrv [-a <address:port>] --sqlitedb [path to DB file]")
		flag.PrintDefaults()
	}

	var addr string
	var db string

	flag.StringVar(&addr, "a", ":8080", "address for server")
	flag.StringVar(&addr, "addr", ":8080", "address for server")
	flag.StringVar(&db, "sqlitedb", "app.db", "path to SQLite database for application data")
	flag.Parse()

	err := web.Start(addr, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "simplesrv: unable to start server: %v\n", err)
		os.Exit(1)
	}
}
