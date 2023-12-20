package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"snippetbox.martinlabs.io/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Add snippets field to application stuct, makes SnippetModel available to our handlers
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Define what port the web server will run on
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", "web:neil@/snippetbox?parseTime=true", "MySQL data source name")
	// Parse the command-line flag
	flag.Parse()

	// Use the slog.New() function to initalise a new structured logger to write to the stream
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()
	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// And add it to the application dependencies.
	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db}, templateCache: templateCache,
	}
	logger.Info("starting server", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// OpenDB function wraps sql.Open() and return sql.DB connection pool for a given dsn
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}
