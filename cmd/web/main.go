package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	// Define what port the web server will run on
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse the command-line flag
	flag.Parse()

	// Use the slog.New() function to initalise a new structured logger to write to the stream
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	// Create a file server which servers files out of the /static/ directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() to register the file server as a handler for all URL paths within "/static/"
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	logger.Info("starting server on", "addr", *addr)

	// Pass the addr pointer to the http.ListenAndServe()
	err := http.ListenAndServe(*addr, mux)

	// Logs any error message returned by http.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
