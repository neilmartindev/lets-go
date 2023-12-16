package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snipper into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// This will return a new snipppet into the database
func (m *SnippetModel) Get(id int) (Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	// Write the sQL statement we want to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute our
	// SQL statement.

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Ensure the sql.Rows resultset is always closed
	defer rows.Close()

	// Initialise an empty slice to hold the snippet structs
	var snippets []Snippet

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and subsequent) row to be acted on by
	// the rows.Scan() method. Once the iteration is done it automatically closes.
	for rows.Next() {
		var s Snippet

		// rows.Scan() copies the values from each field into the new Snppet object
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Appended to the slice
		snippets = append(snippets, s)
	}

	// Retrives any error that was encountered during the iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
