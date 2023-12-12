package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define what port the web server will run on
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse the command-line flag
	flag.Parse()

	mux := http.NewServeMux()

	// Create a file server which servers files out of the /static/ directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() to register the file server as a handler for all URL paths within "/static/"
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on", *addr)

	// Pass the addr pointer to the http.ListenAndServe()
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
